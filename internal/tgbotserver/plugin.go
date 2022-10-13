package tgbotserver

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"go.uber.org/zap"
)

const PluginName = "tgbotserver"

type Plugin struct {
	cfg     *Config
	polling *pollingMode
	webhook *WebhookConfig
}

func (p *Plugin) Init(cfg config.Configurer, bot *tgbotapi.BotAPI, client *asynq.Client, log *zap.Logger) error {
	const op = errors.Op("tgbotserver plugin init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var err error
	if err = cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, errors.Init, err)
	}
	p.cfg.InitDefault()

	switch p.cfg.Mode {
	case webhook:
		if err = cfg.UnmarshalKey(PluginName, &p.webhook); err != nil {
			return errors.E(op, errors.Init, err)
		}
	case polling:
		p.polling = &pollingMode{
			log:    log.Named("polling"),
			bot:    bot,
			client: client,
		}
		if err = cfg.UnmarshalKey(PluginName, &p.polling.cfg); err != nil {
			return errors.E(op, errors.Init, err)
		}
		p.polling.cfg.InitDefault()
	default:
		return errors.E(op, errors.Disabled)
	}

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	if p.cfg.Mode == webhook {
		const op = errors.Op("tgbotserver plugin serve")
		errCh <- errors.E(op, errors.Serve, errors.Str("webhook not implemented"))
	} else {
		p.polling.start()
	}

	return errCh
}

func (p *Plugin) Stop() error {
	if p.cfg.Mode == webhook {
	} else {
		p.polling.stop()
	}
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
