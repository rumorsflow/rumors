package logger

import (
	"context"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
)

const PluginName = "log"

type Plugin struct {
	base     *slog.Logger
	cfg      *logger.Config
	channels logger.ChannelConfig
}

func (p *Plugin) Init(cfg config.Configurer) error {
	const op = errors.Op("logger_plugin_init")
	var err error

	if !cfg.Has(PluginName) {
		p.cfg = &logger.Config{}

		p.base, err = p.cfg.Logger()
		if err != nil {
			return errors.E(op, err)
		}

		slog.SetDefault(p.base)

		return nil
	}

	if err = cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, err)
	}

	if err = cfg.UnmarshalKey(PluginName, &p.channels); err != nil {
		return errors.E(op, err)
	}

	p.base, err = p.cfg.Logger()
	if err != nil {
		return errors.E(op, err)
	}

	slog.SetDefault(p.base)

	return nil
}

func (p *Plugin) Serve() chan error {
	return make(chan error, 1)
}

func (p *Plugin) Stop(context.Context) error {
	if syncer, ok := p.base.Handler().(logger.HandlerSyncer); ok {
		_ = syncer.Sync()
	}
	return nil
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*logger.Logger)(nil), p.ServiceLogger),
	}
}

func (p *Plugin) ServiceLogger() logger.Logger {
	return logger.NewLogger(p.channels, p.base)
}

func (p *Plugin) Name() string {
	return PluginName
}
