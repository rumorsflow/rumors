package task

import "time"

const DefaultQueue = "default"

type ServerConfig struct {
	Concurrency              int            `mapstructure:"concurrency"`
	Queues                   map[string]int `mapstructure:"queues"`
	StrictPriority           bool           `mapstructure:"strict_priority"`
	HealthCheckInterval      time.Duration  `mapstructure:"health_check_interval"`
	DelayedTaskCheckInterval time.Duration  `mapstructure:"delayed_task_check_interval"`
	GroupGracePeriod         time.Duration  `mapstructure:"group_grace_period"`
	GroupMaxDelay            time.Duration  `mapstructure:"group_max_delay"`
	GroupMaxSize             int            `mapstructure:"group_max_size"`
	GracefulTimeout          time.Duration  `mapstructure:"graceful_timeout"`
}

func (cfg *ServerConfig) Init() {
	if cfg.Queues == nil {
		cfg.Queues = make(map[string]int)
	}

	if _, ok := cfg.Queues[DefaultQueue]; !ok {
		cfg.Queues[DefaultQueue] = 1
	}

	if cfg.GracefulTimeout == 0 {
		cfg.GracefulTimeout = 10 * time.Second
	}
}
