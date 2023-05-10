package pubsub

import (
	"context"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
)

const PluginName = "pubsub"

type Plugin struct {
	pub *Publisher
	sub *Subscriber
}

func (p *Plugin) Init(rdbMaker common.RedisMaker, log logger.Logger) error {
	const op = errors.Op("pubsub_plugin_init")

	l := log.NamedLogger(PluginName)

	pub, err := NewPublisher(rdbMaker, l.WithGroup("publisher"))
	if err != nil {
		return errors.E(op, err)
	}

	sub, err := NewSubscriber(rdbMaker)
	if err != nil {
		return errors.E(op, err)
	}

	p.pub = pub
	p.sub = sub

	return nil
}

func (p *Plugin) Serve() chan error {
	return make(chan error, 1)
}

func (p *Plugin) Stop(context.Context) error {
	return errs.Append(p.pub.Close(), p.sub.Close())
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*common.Pub)(nil), p.Pub),
		dep.Bind((*common.Sub)(nil), p.Sub),
	}
}

func (p *Plugin) Pub() common.Pub {
	return p.pub
}

func (p *Plugin) Sub() common.Sub {
	return p.sub
}

func (p *Plugin) Name() string {
	return PluginName
}
