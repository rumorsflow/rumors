package tgbotsender

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"go.uber.org/zap"
	"html/template"
	"sync"
)

const PluginName = "tgbotsender"

type Plugin struct {
	sync.Mutex

	cfg       *Config
	log       *zap.Logger
	botApi    *tgbotapi.BotAPI
	templates *template.Template
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger, botApi *tgbotapi.BotAPI) error {
	const op = errors.Op("tgbotsender plugin init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, errors.Init, err)
	}

	p.log = log
	p.botApi = botApi

	if err := p.initTemplates(); err != nil {
		return errors.E(op, errors.Init, err)
	}

	return nil
}

func (p *Plugin) Serve() chan error {
	go p.SendView(0, ViewAppStart, nil)
	return make(chan error, 1)
}

func (p *Plugin) Stop() error {
	p.SendView(0, ViewAppStop, nil)
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
