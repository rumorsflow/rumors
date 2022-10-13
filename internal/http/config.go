package http

type Config struct {
	Network    string   `mapstructure:"network"`
	Address    string   `mapstructure:"address"`
	Middleware []string `mapstructure:"middleware"`
}

func (cfg *Config) InitDefault() {
	if cfg.Network == "" {
		cfg.Network = "tcp"
	}
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}
}
