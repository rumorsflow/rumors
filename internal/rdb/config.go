package rdb

import "time"

type Config struct {
	Addrs                 []string      `mapstructure:"addrs"`
	ClientName            string        `mapstructure:"client_name"`
	DB                    int           `mapstructure:"db"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	SentinelUsername      string        `mapstructure:"sentinel_username"`
	SentinelPassword      string        `mapstructure:"sentinel_password"`
	MaxRetries            int           `mapstructure:"max_retries"`
	MinRetryBackoff       time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff       time.Duration `mapstructure:"max_retry_backoff"`
	DialTimeout           time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout           time.Duration `mapstructure:"read_timeout"`
	WriteTimeout          time.Duration `mapstructure:"write_timeout"`
	ContextTimeoutEnabled bool          `mapstructure:"context_timeout_enabled"`
	PoolFIFO              bool          `mapstructure:"pool_fifo"`
	PoolSize              int           `mapstructure:"pool_size"`
	PoolTimeout           time.Duration `mapstructure:"pool_timeout"`
	MinIdleConns          int           `mapstructure:"min_idle_conns"`
	MaxIdleConns          int           `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime       time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxLifetime       time.Duration `mapstructure:"conn_max_life_time"`
	MaxRedirects          int           `mapstructure:"max_redirects"`
	ReadOnly              bool          `mapstructure:"read_only"`
	RouteByLatency        bool          `mapstructure:"route_by_latency"`
	RouteRandomly         bool          `mapstructure:"route_randomly"`
	MasterName            string        `mapstructure:"master_name"`
	Ping                  bool          `mapstructure:"ping"`
}
