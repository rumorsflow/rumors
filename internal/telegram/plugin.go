package telegram

import (
	"context"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
)

const (
	PluginName = "telegram"

	sectionPoller = "telegram.poller"
)

type Plugin struct {
	sub    *Subscriber
	poller *Poller
	done   chan struct{}
}

func (p *Plugin) Init(cfg config.Configurer, sub common.Sub, uow common.UnitOfWork, client common.Client, log logger.Logger) error {
	const op = errors.Op("telegram_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var c Config
	if err := cfg.UnmarshalKey(PluginName, &c); err != nil {
		return errors.E(op, err)
	}
	c.Init()

	siteRep, err := uow.Repository((*entity.Site)(nil))
	if err != nil {
		return errors.E(op, err)
	}

	chatRep, err := uow.Repository((*entity.Chat)(nil))
	if err != nil {
		return errors.E(op, err)
	}

	l := log.NamedLogger(PluginName)
	bot := NewBot(&c, l)

	p.sub = NewSubscriber(
		bot,
		sub,
		siteRep.(repository.ReadWriteRepository[*entity.Site]),
		chatRep.(repository.ReadWriteRepository[*entity.Chat]),
		l.WithGroup("subscriber"),
	)

	if cfg.Has(sectionPoller) {
		var pollerCfg PollerConfig
		if err = cfg.UnmarshalKey(sectionPoller, &pollerCfg); err != nil {
			return errors.E(op, err)
		}
		pollerCfg.Init()

		p.poller = NewPoller(&pollerCfg, bot, client, l.WithGroup("poller"))
	}

	return nil
}

func (p *Plugin) Serve() chan error {
	p.done = make(chan struct{})

	p.sub.Run(p.done)

	if p.poller != nil {
		p.poller.Run(p.done)
	}

	return make(chan error, 1)
}

func (p *Plugin) Stop(context.Context) error {
	close(p.done)

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}
