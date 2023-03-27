package telegram

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
