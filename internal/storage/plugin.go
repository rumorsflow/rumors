package storage

import (
	"context"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	mongoext "github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sync"
)

const PluginName = "storage"

type Plugin struct {
	cfg             *Config
	log             *zap.Logger
	roomStorage     *roomStorage
	feedStorage     *feedStorage
	feedItemStorage *feedItemStorage
	userStorage     *userStorage
	done            chan struct{}
	wg              sync.WaitGroup
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger, db *mongo.Database) error {
	const op = errors.Op("storage plugin init")

	if cfg.Has(PluginName) {
		if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
			return errors.E(op, errors.Init, err)
		}
	} else {
		p.cfg = new(Config)
	}

	p.log = log
	p.roomStorage = newRoomStorage(db)
	p.feedStorage = newFeedStorage(db)
	p.feedItemStorage = newFeedItemStorage(db)
	p.userStorage = newUserStorage(db, p.cfg.Admins)
	p.done = make(chan struct{})

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.wg.Add(1)
	go p.serve(errCh)

	return errCh
}

func (p *Plugin) Stop() error {
	close(p.done)
	p.wg.Wait()
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

// Provides declares factory methods.
func (p *Plugin) Provides() []any {
	return []any{
		p.RoomStorage,
		p.FeedStorage,
		p.FeedItemStorage,
		p.UserStorage,
	}
}

func (p *Plugin) RoomStorage() RoomStorage {
	return p.roomStorage
}

func (p *Plugin) FeedStorage() FeedStorage {
	return p.feedStorage
}

func (p *Plugin) FeedItemStorage() FeedItemStorage {
	return p.feedItemStorage
}

func (p *Plugin) UserStorage() UserStorage {
	return p.userStorage
}

func (p *Plugin) serve(errCh chan error) {
	defer p.wg.Done()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const op = errors.Op("storage plugin serve")

	fns := []func(context.Context) error{
		p.roomStorage.indexes,
		p.feedStorage.indexes,
		p.feedItemStorage.indexes,
		p.userStorage.indexes,
	}

	var wg sync.WaitGroup
	wg.Add(4)

	for _, fn := range fns {
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				errCh <- errors.E(op, errors.Serve, err)
			}
		}(fn)
	}

	if p.cfg != nil && len(p.cfg.Admins) > 0 {
		go func(s UserStorage, admins []string) {
			wg.Wait()

			c := mongoext.Criteria{
				Filter: mongoext.Filter{&mongoext.Field{Name: "username", Op: mongoext.OpIn, Value: admins}}.Build(),
				Size:   int64(len(admins)),
			}

			if users, err := s.Find(ctx, c); err == nil {
				for _, user := range users {
					if !user.IsGranted(models.AdminRole) {
						user := user
						user.DeleteRoles = user.Roles
						user.Roles = []models.Role{models.AdminRole}
						_ = s.Save(ctx, &user)
					}
				}
			}
		}(p.userStorage, p.cfg.Admins)
	}

	select {
	case <-p.done:
		p.log.Debug("storage done")
	case <-ctx.Done():
		if ctx.Err() != nil {
			errCh <- errors.E(op, errors.Serve, ctx.Err())
		}
	}
}
