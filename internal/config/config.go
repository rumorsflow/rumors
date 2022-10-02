package config

import "time"

type Config struct {
	Debug    bool           `mapstructure:"debug"`
	Log      LogConfig      `mapstructure:"log"`
	Asynq    AsynqConfig    `mapstructure:"asynq"`
	Server   ServerConfig   `mapstructure:"server"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Telegram TelegramConfig `mapstructure:"telegram"`
}

type LogConfig struct {
	Level   string `mapstructure:"level"`
	Console bool   `mapstructure:"console"`
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
	Group       struct {
		Max struct {
			Delay time.Duration `mapstructure:"delay"`
			Size  int           `mapstructure:"size"`
		} `mapstructure:"max"`
		Grace struct {
			Period time.Duration `mapstructure:"period"`
		} `mapstructure:"grace"`
	} `mapstructure:"group"`
}

type AsynqSchedulerConfig struct {
	FeedImporter string `mapstructure:"feed"`
}

type ServerConfig struct {
	Network string `mapstructure:"network"`
	Address string `mapstructure:"address"`
}

type MongoDBConfig struct {
	URI string `mapstructure:"uri"`
}

type TelegramConfig struct {
	Owner int64  `mapstructure:"owner"`
	Token string `mapstructure:"token"`
}
