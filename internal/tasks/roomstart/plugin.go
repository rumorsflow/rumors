package roomstart

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/services/room"
	"go.uber.org/zap"
)

const PluginName = consts.TaskRoomStart

type Plugin struct {
	log     *zap.Logger
	service room.Service
}

func (p *Plugin) Init(log *zap.Logger, service room.Service) error {
	p.log = log
	p.service = service
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message
	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}
	return p.service.ChatMemberUpdated(ctx, *message.Chat, false)
}
