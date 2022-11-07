package token

import "time"

type Config struct {
	PrivateKey string `mapstructure:"private_key"`
	TTL        struct {
		JWT     time.Duration `mapstructure:"jwt"`
		Refresh time.Duration `mapstructure:"refresh"`
	} `mapstructure:"ttl"`
}

func (cfg *Config) InitDefault() {
	if cfg.TTL.JWT == 0 {
		cfg.TTL.JWT = 5 * time.Minute
	}

	if cfg.TTL.Refresh == 0 {
		cfg.TTL.JWT = 7200 * time.Hour // 5 Days
	}
}
