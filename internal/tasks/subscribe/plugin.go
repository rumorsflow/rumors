package subscribe

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/services/room"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const PluginName = consts.TaskSubscribe

type Plugin struct {
	log     *zap.Logger
	service room.Service
	sender  tgbotsender.TelegramSender
}

func (p *Plugin) Init(log *zap.Logger, service room.Service, sender tgbotsender.TelegramSender) error {
	p.log = log
	p.service = service
	p.sender = sender
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

	if message.CommandArguments() == "" {
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, consts.ErrMsgRequiredSource)
		return nil
	}

	if !lo.Contains([]models.RoomType{models.Private, models.Channel}, models.RoomType(message.Chat.Type)) &&
		(message.From == nil || message.From.ID != p.sender.Owner()) {
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewForbidden, nil)
		return nil
	}

	if err := p.service.AddToBroadcastByHost(ctx, message.Chat.ID, message.CommandArguments()); err != nil {
		if errors2.Is(err, room.ErrNotFoundFeeds) {
			p.sender.SendView(message.Chat.ID, tgbotsender.ViewNotFound, fmt.Sprintf(consts.ErrMsgNotFoundSource, message.CommandArguments()))
		} else {
			p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, nil)
		}
		return nil
	}

	p.sender.SendView(message.Chat.ID, tgbotsender.ViewSuccess, consts.SuccessMsgSubscribed)

	return nil
}
