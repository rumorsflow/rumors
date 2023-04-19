package jwt

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
)

type (
	ConfigKey struct{}
	SignerKey struct{}
)

func GetConfig(ctx context.Context, c ...di.Container) (*Config, error) {
	return di.Get[*Config](ctx, ConfigKey{}, c...)
}

func GetSigner(ctx context.Context, c ...di.Container) (Signer, error) {
	return di.Get[Signer](ctx, SignerKey{}, c...)
}

func ConfigActivator(configKey string) *di.Activator {
	return &di.Activator{
		Key: ConfigKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*Config](c.Configurer(), configKey)
			if err != nil {
				return nil, nil, fmt.Errorf("%s error: %w", di.OpFactory, err)
			}

			return cfg, nil, nil
		}),
	}
}

func SignerActivator() *di.Activator {
	return &di.Activator{
		Key: SignerKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := GetConfig(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return NewSigner(cfg.GetPrivateKey()), nil, nil
		}),
	}
}
