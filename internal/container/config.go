package container

import (
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
	"time"
)

type Config struct {
	GracePeriod time.Duration
	PrintGraph  bool
	LogLevel    slog.Leveler
}

const (
	endureKey          = "endure"
	defaultGracePeriod = time.Second * 30
)

// NewConfig creates endure container configuration.
func NewConfig(cfgFile, prefix string) (*Config, error) {
	cfg, err := config.NewConfigurer(
		"",
		0,
		config.WithPath(cfgFile),
		config.WithPrefix(prefix),
	)
	if err != nil {
		return nil, err
	}

	if !cfg.Has(endureKey) {
		return &Config{ // return config with defaults
			GracePeriod: defaultGracePeriod,
			PrintGraph:  false,
			LogLevel:    slog.LevelError,
		}, nil
	}

	rrCfgEndure := struct {
		GracePeriod time.Duration `mapstructure:"grace_period"`
		PrintGraph  bool          `mapstructure:"print_graph"`
		LogLevel    string        `mapstructure:"log_level"`
	}{}

	if err = cfg.UnmarshalKey(endureKey, &rrCfgEndure); err != nil {
		return nil, err
	}

	if rrCfgEndure.GracePeriod == 0 {
		rrCfgEndure.GracePeriod = defaultGracePeriod
	}

	if rrCfgEndure.LogLevel == "" {
		rrCfgEndure.LogLevel = "error"
	}

	return &Config{
		GracePeriod: rrCfgEndure.GracePeriod,
		PrintGraph:  rrCfgEndure.PrintGraph,
		LogLevel:    logger.ToLeveler(rrCfgEndure.LogLevel),
	}, nil
}
