package config

type Config struct {
	Debug    bool           `mapstructure:"debug"`
	DB       DBConfig       `mapstructure:"db"`
	Log      LogConfig      `mapstructure:"log"`
	Asynq    AsynqConfig    `mapstructure:"asynq"`
	Server   ServerConfig   `mapstructure:"server"`
	Telegram TelegramConfig `mapstructure:"telegram"`
}

type DBConfig struct {
	Path string `mapstructure:"path"`
	Main string `mapstructure:"main"`
	Data string `mapstructure:"data"`
}

type LogConfig struct {
	Level   string `mapstructure:"level"`
	Colored bool   `mapstructure:"colored"`
}

type AsynqConfig struct {
	Redis     RedisConfig          `mapstructure:"redis"`
	Server    AsynqServerConfig    `mapstructure:"server"`
	Scheduler AsynqSchedulerConfig `mapstructure:"scheduler"`
}

type RedisConfig struct {
	Network  string `mapstructure:"network"`
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AsynqServerConfig struct {
	Concurrency int `mapstructure:"concurrency"`
}

type AsynqSchedulerConfig struct {
	Cronspec string `mapstructure:"cron"`
	TaskName string `mapstructure:"name"`
}

type ServerConfig struct {
	Network string `mapstructure:"network"`
	Address string `mapstructure:"address"`
}

type TelegramConfig struct {
	Owner int64  `mapstructure:"owner"`
	Token string `mapstructure:"token"`
}
