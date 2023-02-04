package di

import (
	"context"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"sync"
)

var (
	ioc  Container
	once sync.Once
)

func Init(cfg config.Configurer) {
	once.Do(func() {
		ioc = NewContainer(cfg)
	})
}

func Default() Container {
	if ioc == nil {
		panic("default container is nil, call di.Init")
	}
	return ioc
}

func Configurer() config.Configurer {
	return Default().Configurer()
}

func Activators(activators ...*Activator) error {
	return Default().Activators(activators...)
}

func Register(key any, factory Factory) error {
	return Default().Register(key, factory)
}

func Close(ctx context.Context) error {
	return Default().Close(ctx)
}

func New[T any](ctx context.Context, key any, c ...Container) (value T, err error) {
	i := Default()
	if len(c) > 0 {
		i = c[0]
	}
	v, err := i.New(ctx, key)
	if err != nil {
		return value, err
	}
	return v.(T), nil
}

func MustNew[T any](ctx context.Context, key any, c ...Container) T {
	return Must(New[T](ctx, key, c...))
}

func Get[T any](ctx context.Context, key any, c ...Container) (value T, err error) {
	i := Default()
	if len(c) > 0 {
		i = c[0]
	}
	v, err := i.Get(ctx, key)
	if err != nil {
		return value, err
	}
	return v.(T), nil
}

func MustGet[T any](ctx context.Context, key any, c ...Container) T {
	return Must(Get[T](ctx, key, c...))
}

func Has(key any, c ...Container) bool {
	if len(c) == 0 {
		return Default().Has(key)
	}
	return c[0].Has(key)
}
