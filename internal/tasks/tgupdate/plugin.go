package tgupdate

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"go.uber.org/zap"
)

const PluginName = consts.TaskTelegramUpdate

var commands = map[string]string{
	consts.TgCmdStart:       consts.TaskRoomStart,
	consts.TgCmdRumors:      consts.TaskRumors,
	consts.TgCmdSources:     consts.TaskSources,
	consts.TgCmdSubscribed:  consts.TaskSubscribed,
	consts.TgCmdSubscribe:   consts.TaskSubscribe,
	consts.TgCmdUnsubscribe: consts.TaskUnsubscribe,
}

type Plugin struct {
	log    *zap.Logger
	client *asynq.Client
}

func (p *Plugin) Init(log *zap.Logger, client *asynq.Client) error {
	p.log = log
	p.client = client
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var update tgbotapi.Update
	if err := json.Unmarshal(task.Payload(), &update); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}

	if update.Message != nil {
		return p.message(ctx, update.Message)
	}

	if update.EditedMessage != nil {
		return p.message(ctx, update.EditedMessage)
	}

	if update.ChannelPost != nil {
		return p.message(ctx, update.ChannelPost)
	}

	if update.EditedChannelPost != nil {
		return p.message(ctx, update.EditedChannelPost)
	}

	if update.MyChatMember != nil {
		return p.chatMemberUpdated(ctx, update.MyChatMember)
	}

	if update.ChatMember != nil {
		return p.chatMemberUpdated(ctx, update.ChatMember)
	}

	return nil
}

func (p *Plugin) message(ctx context.Context, message *tgbotapi.Message) error {
	if !message.IsCommand() {
		return nil
	}
	if cmd, ok := commands[message.Command()]; ok {
		return p.enqueue(ctx, cmd, message)
	}
	return nil
}

func (p *Plugin) chatMemberUpdated(ctx context.Context, member *tgbotapi.ChatMemberUpdated) error {
	return p.enqueue(ctx, consts.TaskRoomUpdated, member)
}

func (p *Plugin) enqueue(ctx context.Context, taskName string, data any) (err error) {
	payload, _ := json.Marshal(data)
	_, err = p.client.EnqueueContext(ctx, asynq.NewTask(taskName, payload))
	if err != nil {
		p.log.Error(
			"error due to enqueue task",
			zap.Error(err),
			zap.String("new_task", taskName),
			zap.ByteString("new_payload", payload),
		)
	}
	return
}
