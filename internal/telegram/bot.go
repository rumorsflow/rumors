package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
	"html/template"
	"net/http"
	"time"
)

type Bot struct {
	cfg       *Config
	logger    *slog.Logger
	api       *tgbotapi.BotAPI
	templates *template.Template
}

func NewBot(cfg *Config) *Bot {
	cfg.Init()

	log := logger.WithGroup("telegram").WithGroup("bot")
	_ = tgbotapi.SetLogger(&telegramLogger{logger: log})

	api := &tgbotapi.BotAPI{
		Token:  cfg.Token,
		Debug:  logger.IsDebug(),
		Client: &http.Client{},
		Buffer: 100,
	}
	api.SetAPIEndpoint(tgbotapi.APIEndpoint)

	self, err := api.GetMe()
	if err != nil {
		panic(fmt.Errorf("%s error: %w", OpBotNew, err))
	}
	api.Self = self

	return &Bot{
		cfg:    cfg,
		logger: log,
		api:    api,
	}
}

func (b *Bot) OwnerID() int64 {
	return b.cfg.OwnerID
}

func (b *Bot) Me() tgbotapi.User {
	return b.api.Self
}

func (b *Bot) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return b.api.Request(c)
}

func (b *Bot) Send(message Message) error {
	messages, err := message.chattable(b)
	if err != nil {
		return err
	}
	for _, msg := range messages {
		if err = b.Chattable(msg); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) Chattable(c tgbotapi.Chattable) error {
	b.logger.Debug("bot send message", "message", c)
	return b.chattable(c, 0)
}

func (b *Bot) chattable(c tgbotapi.Chattable, retry uint) error {
	if _, err := b.Request(c); err != nil {
		var res []byte
		if e, ok := err.(*tgbotapi.Error); ok {
			if e.RetryAfter > 0 && retry < b.cfg.Retry {
				time.Sleep(time.Duration(e.RetryAfter+1) * time.Second)

				return b.chattable(c, retry+1)
			}
			res, _ = json.Marshal(e)
		}

		b.logger.Error("bot request error", err, "message", c, "res_err", string(res))

		return fmt.Errorf("%s error: %w", OpBotSend, err)
	}

	return nil
}
