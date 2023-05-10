package config

import (
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"time"
)

const PluginName string = "config"

type Plugin struct {
	Path      string
	Prefix    string
	Type      string
	ReadInCfg []byte
	// user defined Flags in the form of <option>.<key> = <value>
	// which overwrites initial config key
	Flags []string

	// Timeout ...
	Timeout time.Duration
	Version string

	cfg config.Configurer
}

// Init config provider.
func (p *Plugin) Init() error {
	const op = errors.Op("config_plugin_init")

	if cfg, err := config.NewConfigurer(
		p.Version,
		p.Timeout,
		config.WithPath(p.Path),
		config.WithPrefix(p.Prefix),
		config.WithConfigType(p.Type),
		config.WithReadInCfg(p.ReadInCfg),
		config.WithFlags(p.Flags),
	); err != nil {
		return errors.E(op, err)
	} else {
		p.cfg = cfg
	}

	return nil
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*config.Configurer)(nil), p.ServiceConfigurer),
	}
}

func (p *Plugin) ServiceConfigurer() config.Configurer {
	return p.cfg
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
