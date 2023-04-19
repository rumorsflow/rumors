package telegram

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
)

const ConfigKey = "telegram"

type (
	BotKey struct{}
	SubKey struct{}
)

func GetBot(ctx context.Context, c ...di.Container) (*Bot, error) {
	return di.Get[*Bot](ctx, BotKey{}, c...)
}

func GetSub(ctx context.Context, c ...di.Container) (*Subscriber, error) {
	return di.Get[*Subscriber](ctx, SubKey{}, c...)
}

func BotActivator() *di.Activator {
	return &di.Activator{
		Key: BotKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*Config](c.Configurer(), ConfigKey)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			return NewBot(cfg), nil, nil
		}),
	}
}

func SubActivator() *di.Activator {
	return &di.Activator{
		Key: SubKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			bot, err := GetBot(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			sub, err := pubsub.GetSub(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			siteRepo, err := db.GetSiteRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			chatRepo, err := db.GetChatRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return NewSubscriber(bot, sub, siteRepo, chatRepo), nil, nil
		}),
	}
}
