package jwt

import (
	"crypto/rsa"
	"time"
)

type Config struct {
	PrivateKey      string        `mapstructure:"private_key"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`

	privateKey *rsa.PrivateKey
}

func (cfg *Config) Init() {
	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = 1 * time.Minute
	}

	if cfg.RefreshTokenTTL == 0 {
		cfg.RefreshTokenTTL = 24 * time.Hour
	}
}

func (cfg *Config) GetPrivateKey() *rsa.PrivateKey {
	if cfg.privateKey == nil {
		privateKey, err := GetRSAPrivateKey(cfg.PrivateKey)
		if err != nil {
			panic(err)
		}
		cfg.privateKey = privateKey
	}

	return cfg.privateKey
}

func (cfg *Config) GetPublicKey() *rsa.PublicKey {
	return &cfg.GetPrivateKey().PublicKey
}
