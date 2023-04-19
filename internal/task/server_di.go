package task

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
)

const ConfigServerKey = "task.server"

type ServerKey struct{}

func GetServer(ctx context.Context, c ...di.Container) (*Server, error) {
	return di.Get[*Server](ctx, ServerKey{}, c...)
}

func ServerActivator() *di.Activator {
	return &di.Activator{
		Key: ServerKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*ServerConfig](c.Configurer(), ConfigServerKey)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return NewServer(cfg, rdbMaker), nil, nil
		}),
	}
}
