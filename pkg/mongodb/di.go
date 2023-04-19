package mongodb

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
)

type DatabaseKey struct{}

func New(ctx context.Context, c ...di.Container) (*Database, error) {
	return di.New[*Database](ctx, DatabaseKey{}, c...)
}

func Get(ctx context.Context, c ...di.Container) (*Database, error) {
	return di.Get[*Database](ctx, DatabaseKey{}, c...)
}

func Activator(configKey string) *di.Activator {
	return &di.Activator{
		Key: DatabaseKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*Config](c.Configurer(), configKey)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			db, err := NewDatabase(ctx, cfg)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			return db, db, nil
		}),
	}
}
