package tgbotserver

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"go.uber.org/zap"
	"time"
)

type pollingMode struct {
	cfg    *PollingConfig
	bot    *tgbotapi.BotAPI
	client *asynq.Client
	log    *zap.Logger
	done   chan struct{}
}

func (p *pollingMode) start() {
	p.done = make(chan struct{})
	go p.getUpdates(p.cfg.BuildUpdateConfig())
}

func (p *pollingMode) stop() {
	p.log.Info("Stopping the update receiver routine...")
	close(p.done)
}

func (p *pollingMode) getUpdates(config tgbotapi.UpdateConfig) {
	for {
		select {
		case <-p.done:
			return
		default:
		}

		updates, err := p.bot.GetUpdates(config)
		if err != nil {
			p.log.Error("Failed to get updates, retrying in 3 seconds...", zap.Error(err))
			time.Sleep(time.Second * 3)
			continue
		}

		for _, update := range updates {
			if update.UpdateID >= config.Offset {
				config.Offset = update.UpdateID + 1

				payload, _ := json.Marshal(update)
				taskId := asynq.TaskID(fmt.Sprintf("%s:%d", consts.TaskTelegramUpdate, update.UpdateID))

				p.log.Debug("telegram update", zap.ByteString("update", payload))

				if _, err = p.client.Enqueue(asynq.NewTask(consts.TaskTelegramUpdate, payload), taskId); err != nil {
					p.log.Error("error due to enqueue update", zap.Error(err))
				}
			}
		}
	}
}
