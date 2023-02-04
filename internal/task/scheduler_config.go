package task

import "time"

type SchedulerConfig struct {
	SyncInterval time.Duration `mapstructure:"sync_interval"`
}

func (cfg *SchedulerConfig) Init() {
	if cfg.SyncInterval == 0 {
		cfg.SyncInterval = 5 * time.Minute
	}
}
