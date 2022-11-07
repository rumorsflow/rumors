package token

import (
	"github.com/go-redis/redis/v8"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"github.com/rumorsflow/rumors/internal/pkg/jwt"
	"github.com/rumorsflow/rumors/internal/storage"
)

const PluginName = "token"

type Plugin struct {
	cfg         *Config
	userStorage storage.UserStorage
	storage     Storage
	signer      jwt.Signer
	verifier    jwt.Verifier
}

func (p *Plugin) Init(cfg config.Configurer, s storage.UserStorage, client redis.UniversalClient) error {
	const op = errors.Op("token service plugin init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, errors.Init, err)
	}
	p.cfg.InitDefault()

	prvKey, err := jwt.GetRSAPrivateKey(p.cfg.PrivateKey)
	if err != nil {
		return errors.E(op, errors.Init, err)
	}

	p.signer, err = jwt.NewSigner(prvKey)
	if err != nil {
		return errors.E(op, errors.Init, err)
	}

	p.verifier, err = jwt.NewVerifier(&prvKey.PublicKey)
	if err != nil {
		return errors.E(op, errors.Init, err)
	}

	p.userStorage = s
	p.storage = &tokenStorage{client: client}

	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
