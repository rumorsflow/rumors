package telegram

import "time"

type Config struct {
	Token   string `mapstructure:"token"`
	OwnerID int64  `mapstructure:"owner"`
	Retry   uint   `mapstructure:"retry"`
}

func (cfg *Config) Init() {
	if cfg.Token == "" {
		panic("telegram bot API is required")
	}

	if cfg.OwnerID == 0 {
		panic("telegram owner ID is required")
	}

	if cfg.Retry == 0 {
		cfg.Retry = 3
	}
}

type PollerConfig struct {
	OnlyOwner      bool          `mapstructure:"only_owner"`
	Buffer         int           `mapstructure:"buffer"`
	Limit          int           `mapstructure:"limit"`
	Timeout        time.Duration `mapstructure:"timeout"`
	AllowedUpdates []string      `mapstructure:"allowed_updates"`
}

func (cfg *PollerConfig) Init() {
	if cfg.Limit < 1 || cfg.Limit > 100 {
		cfg.Limit = 100
	}

	if cfg.Buffer < 0 {
		cfg.Buffer = 100
	}
}
