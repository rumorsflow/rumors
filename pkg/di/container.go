package di

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"go.uber.org/multierr"
	"golang.org/x/exp/slog"
	"reflect"
	"sync"
)

const (
	OpActivators errs.Op = "container: activators"
	OpRegister   errs.Op = "container: register"
	OpFactory    errs.Op = "container: factory"
	OpNew        errs.Op = "container: new"
	OpGet        errs.Op = "container: get"
	OpClose      errs.Op = "container: close"
)

type Closer interface {
	Close(ctx context.Context) error
}

type CloserFunc func(ctx context.Context) error

func (c CloserFunc) Close(ctx context.Context) error {
	return c(ctx)
}

type Factory interface {
	New(context.Context, Container) (any, Closer, error)
}

type FactoryFunc func(ctx context.Context, c Container) (any, Closer, error)

func (f FactoryFunc) New(ctx context.Context, c Container) (any, Closer, error) {
	return f(ctx, c)
}

type Activator struct {
	Key     any
	Factory Factory
}

type Container interface {
	Configurer() config.Configurer
	Activators(activators ...*Activator) error
	Register(key any, factory Factory) error
	Close(ctx context.Context) error
	New(ctx context.Context, key any) (any, error)
	Get(ctx context.Context, key any) (any, error)
	Has(key any) bool
}

type closerEntry struct {
	key    any
	closer Closer
}

func newCloser(key any, closer Closer) *closerEntry {
	return &closerEntry{key: key, closer: closer}
}

type container struct {
	configurer config.Configurer
	logger     *slog.Logger
	factories  map[any]Factory
	values     map[any]any
	closers    []*closerEntry
	mu         sync.RWMutex
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func ToNilCloser[T any](value T, err error) (T, Closer, error) {
	return value, nil, err
}

func NewContainer(cfg config.Configurer) Container {
	return &container{
		configurer: cfg,
		logger:     logger.WithGroup("di").WithGroup("container"),
		factories:  map[any]Factory{},
		values:     map[any]any{},
	}
}

func (c *container) Configurer() config.Configurer {
	return c.configurer
}

func (c *container) Activators(activators ...*Activator) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, activator := range activators {
		if activator == nil {
			return errs.Errorf(OpActivators, "Activator `%d` is nil", i)
		}
		if err := c.register(activator.Key, activator.Factory); err != nil {
			return errs.E(OpActivators, err)
		}
	}
	return nil
}

func (c *container) Register(key any, factory Factory) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.register(key, factory); err != nil {
		return err
	}
	return nil
}

func (c *container) Close(ctx context.Context) (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.info("closing")

	for i := len(c.closers) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			err = multierr.Append(err, ctx.Err())
			break
		default:
		}

		err = multierr.Append(err, c.closers[i].closer.Close(ctx))

		c.info("`%s` Close called", str(c.closers[i].key))
	}

	c.closers = nil

	if err != nil {
		err = errs.E(OpClose, err)
	}
	return
}

func (c *container) New(ctx context.Context, key any) (value any, err error) {
	if err = checkKey(key); err != nil {
		return nil, errs.E(OpNew, err)
	}

	var closer Closer
	if value, closer, err = c.new(ctx, key); err == nil && closer != nil {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.closers = append(c.closers, newCloser(key, closer))
	}

	return
}

func (c *container) Get(ctx context.Context, key any) (any, error) {
	if err := checkKey(key); err != nil {
		return nil, errs.E(OpGet, err)
	}

	c.mu.RLock()
	if value, ok := c.values[key]; ok {
		c.mu.RUnlock()

		c.debug("value for `%s` was found", str(key))

		return value, nil
	}
	c.mu.RUnlock()

	value, closer, err := c.new(ctx, key)
	if err != nil {
		return nil, errs.E(OpGet, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value
	if closer != nil {
		c.closers = append(c.closers, newCloser(key, closer))
	}

	return value, nil
}

func (c *container) Has(key any) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if checkKey(key) != nil {
		return false
	}

	_, ok := c.factories[key]
	return ok
}

func (c *container) register(key any, factory Factory) error {
	if err := checkKey(key); err != nil {
		return errs.E(OpRegister, err)
	}

	if factory == nil {
		return errs.Errorf(OpRegister, "factory `%s` is nil", str(key))
	}

	if _, ok := c.factories[key]; ok {
		return errs.Errorf(OpRegister, "factory `%s` already exists", str(key))
	}

	c.factories[key] = factory

	c.info("factory `%s` registered", str(key))

	return nil
}

func (c *container) new(ctx context.Context, key any) (any, Closer, error) {
	c.mu.RLock()
	factory, ok := c.factories[key]
	c.mu.RUnlock()

	if ok {
		value, closer, err := factory.New(ctx, c)

		c.debug("factory `%s` called", str(key))

		if err != nil {
			return nil, nil, errs.E(errs.Errorf(OpNew, "factory `%s`", str(key)), err)
		}
		return value, closer, nil
	}

	return nil, nil, errs.Errorf(OpNew, "factory `%s` not found", str(key))
}

func (c *container) info(format string, a ...any) {
	c.logger.Info(fmt.Sprintf(format, a...))
}

func (c *container) debug(format string, a ...any) {
	c.logger.Debug(fmt.Sprintf(format, a...))
}

func checkKey(key any) error {
	if key == nil || !reflect.TypeOf(key).Comparable() {
		return errors.Errorf("key `%s` is not comparable", str(key))
	}
	return nil
}

func str(k any) string {
	if s, ok := k.(string); ok {
		return s
	}
	return reflect.TypeOf(k).String()
}
