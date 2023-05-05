package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
)

type (
	UniversalClientKey struct{}
	MakerKey           struct{}
)

func NewUniversalClient(ctx context.Context, c ...di.Container) (redis.UniversalClient, error) {
	return di.New[redis.UniversalClient](ctx, UniversalClientKey{}, c...)
}

func GetUniversalClient(ctx context.Context, c ...di.Container) (redis.UniversalClient, error) {
	return di.Get[redis.UniversalClient](ctx, UniversalClientKey{}, c...)
}

func GetMaker(ctx context.Context, c ...di.Container) (*UniversalClientMaker, error) {
	return di.Get[*UniversalClientMaker](ctx, MakerKey{}, c...)
}

func UniversalClientActivator(configKey string) *di.Activator {
	return &di.Activator{
		Key: UniversalClientKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*Config](c.Configurer(), configKey)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			client := New(cfg)

			if cfg.Ping {
				if res, err := client.Ping(ctx).Result(); err != nil || res != "PONG" {
					return nil, nil, fmt.Errorf("%s could not check Redis server. error: %w", di.OpFactory, err)
				}
			}

			return client, nil, nil
		}),
	}
}

func MakerActivator() *di.Activator {
	return &di.Activator{
		Key: MakerKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			return &UniversalClientMaker{c: c}, nil, nil
		}),
	}
}
