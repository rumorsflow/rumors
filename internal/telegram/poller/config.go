package poller

import (
	"time"
)

type Config struct {
	OnlyOwner      bool          `mapstructure:"only_owner"`
	Buffer         int           `mapstructure:"buffer"`
	Limit          int           `mapstructure:"limit"`
	Timeout        time.Duration `mapstructure:"timeout"`
	AllowedUpdates []string      `mapstructure:"allowed_updates"`
}

func (cfg *Config) Init() {
	if cfg.Limit < 1 || cfg.Limit > 100 {
		cfg.Limit = 100
	}

	if cfg.Buffer < 0 {
		cfg.Buffer = 100
	}
}
