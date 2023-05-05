package poller

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
	"time"
)

type TelegramPoller struct {
	cfg    *Config
	bot    *telegram.Bot
	client *task.Client
	logger *slog.Logger
	update chan tgbotapi.Update
	done   chan struct{}
}

func NewTelegramPoller(cfg *Config, bot *telegram.Bot, client *task.Client) *TelegramPoller {
	cfg.Init()

	p := &TelegramPoller{
		cfg:    cfg,
		bot:    bot,
		client: client,
		logger: logger.WithGroup("telegram").WithGroup("poller"),
		done:   make(chan struct{}),
	}

	if cfg.Buffer == 0 {
		p.update = make(chan tgbotapi.Update)
	} else {
		p.update = make(chan tgbotapi.Update, cfg.Buffer)
	}

	return p
}

func (p *TelegramPoller) Poll(ctx context.Context) error {
	go p.listen()

	config := tgbotapi.UpdateConfig{
		Limit:          p.cfg.Limit,
		Timeout:        int(p.cfg.Timeout / time.Second),
		AllowedUpdates: p.cfg.AllowedUpdates,
	}

	defer p.logger.Info("telegram poller stopped")

	for {
		select {
		case <-ctx.Done():
			close(p.done)
			return nil
		default:
		}

		updates, err := p.getUpdates(config)
		if err != nil {
			var retry time.Duration = 3
			if tge, ok := err.(*tgbotapi.Error); ok && tge.RetryAfter > 0 {
				retry = time.Duration(tge.RetryAfter)
			}

			p.logger.Error(fmt.Sprintf("failed to get updates, retrying in %d seconds...", retry), "err", err)
			time.Sleep(retry * time.Second)

			continue
		}

		for _, update := range updates {
			if update.UpdateID >= config.Offset {
				config.Offset = update.UpdateID + 1
				p.update <- update
			}
		}
	}
}

func (p *TelegramPoller) getUpdates(config tgbotapi.UpdateConfig) ([]tgbotapi.Update, error) {
	resp, err := p.bot.Request(config)
	if err != nil {
		return nil, err
	}

	var updates []tgbotapi.Update
	err = json.Unmarshal(resp.Result, &updates)

	return updates, err
}

func (p *TelegramPoller) listen() {
	ctx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-p.done:
			cancel()

			if err := ctx.Err(); err != nil && !errs.IsCanceledOrDeadline(err) {
				p.logger.Error("failed to enqueue update", "err", ctx.Err())
			}

			return
		case update := <-p.update:
			if update.Message != nil {
				p.message(ctx, update.Message, update.UpdateID)
			} else if update.EditedMessage != nil {
				p.message(ctx, update.EditedMessage, update.UpdateID)
			} else if update.ChannelPost != nil {
				p.message(ctx, update.ChannelPost, update.UpdateID)
			} else if update.EditedChannelPost != nil {
				p.message(ctx, update.EditedChannelPost, update.UpdateID)
			} else if update.MyChatMember != nil {
				p.client.EnqueueTgMemberEdit(ctx, update.MyChatMember, update.UpdateID)
			} else if update.ChatMember != nil {
				p.client.EnqueueTgMemberEdit(ctx, update.ChatMember, update.UpdateID)
			}
		}
	}
}

func (p *TelegramPoller) message(ctx context.Context, message *tgbotapi.Message, updateID int) {
	if message == nil || message.Chat == nil || !message.IsCommand() {
		return
	}

	switch message.Command() {
	case task.TgCmdStart:
		p.client.EnqueueTgMemberNew(ctx, message.Chat, updateID)
	case task.TgCmdRumors, task.TgCmdSites, task.TgCmdSub, task.TgCmdOn, task.TgCmdOff:
		if p.cfg.OnlyOwner && !(message.Chat.ID == p.bot.OwnerID() || (message.From != nil && message.From.ID == p.bot.OwnerID())) {
			p.logger.Warn("access denied", "message", message)
			return
		}

		p.client.EnqueueTgCmd(ctx, message, updateID)
	}
}
