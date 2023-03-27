package poller

import (
	"context"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
)

const ConfigKey = "telegram.poller"

type Key struct{}

func GetTelegramPoller(ctx context.Context, c ...di.Container) (*TelegramPoller, error) {
	return di.Get[*TelegramPoller](ctx, Key{}, c...)
}

func TelegramPollerActivator() *di.Activator {
	return &di.Activator{
		Key: Key{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*Config](c.Configurer(), ConfigKey)
			if err != nil {
				return nil, nil, errs.E(di.OpFactory, err)
			}

			bot, err := telegram.GetBot(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			client, err := task.GetClient(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return NewTelegramPoller(cfg, bot, client), nil, nil
		}),
	}
}
