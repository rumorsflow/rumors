package storage

import (
	"context"
	"github.com/roadrunner-server/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"sync"
)

const PluginName = "app_storage"

type Plugin struct {
	log             *zap.Logger
	roomStorage     *roomStorage
	feedStorage     *feedStorage
	feedItemStorage *feedItemStorage
	userStorage     *userStorage
	done            chan struct{}
	wg              sync.WaitGroup
}

func (p *Plugin) Init(log *zap.Logger, db *mongo.Database) error {
	p.log = log
	p.roomStorage = newRoomStorage(db)
	p.feedStorage = newFeedStorage(db)
	p.feedItemStorage = newFeedItemStorage(db)
	p.userStorage = newUserStorage(db)
	p.done = make(chan struct{})

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	p.wg.Add(1)
	go p.indexes(errCh)

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

func (p *Plugin) indexes(errCh chan error) {
	defer p.wg.Done()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const op = errors.Op("app storage plugin indexes")

	fns := []func(context.Context) error{
		p.roomStorage.indexes,
		p.feedStorage.indexes,
		p.feedItemStorage.indexes,
		p.userStorage.indexes,
	}

	for _, fn := range fns {
		go func(fn func(context.Context) error) {
			if err := fn(ctx); err != nil {
				errCh <- errors.E(op, errors.Serve, err)
			}
		}(fn)
	}

	select {
	case <-p.done:
		p.log.Debug("app storage done")
	case <-ctx.Done():
		if ctx.Err() != nil {
			errCh <- errors.E(op, errors.Serve, ctx.Err())
		}
	}
}
