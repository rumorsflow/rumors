package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/rs/zerolog/log"
)

var commands = map[string]string{
	consts.TgCmdAdd:    consts.TaskFeedAdd,
	consts.TgCmdRoom:   consts.TaskRoomView,
	consts.TgCmdFeed:   consts.TaskFeedView,
	consts.TgCmdRumors: consts.TaskFeedItemView,
	consts.TgCmdStart:  consts.TaskRoomAdd,
	"creator":          consts.TaskRoomAdd,
	"member":           consts.TaskRoomAdd,
	"administrator":    consts.TaskRoomAdd,
	"restricted":       consts.TaskRoomAdd,
	"left":             consts.TaskRoomLeft,
	"kicked":           consts.TaskRoomLeft,
}

type TelegramUpdateHandler struct {
	Client  *asynq.Client
	Emitter emitter.Emitter
	Owner   int64
}

func (h *TelegramUpdateHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var update tgbotapi.Update
	if err := json.Unmarshal(task.Payload(), &update); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	} else {
		l.Info().Msg("update received")
	}

	if update.Message != nil {
		return h.message(ctx, update.Message)
	}

	if update.EditedMessage != nil {
		return h.message(ctx, update.EditedMessage)
	}

	if update.ChannelPost != nil {
		return h.message(ctx, update.ChannelPost)
	}

	if update.EditedChannelPost != nil {
		return h.message(ctx, update.EditedChannelPost)
	}

	if update.MyChatMember != nil {
		return h.chatMemberUpdated(ctx, update.MyChatMember)
	}

	if update.ChatMember != nil {
		return h.chatMemberUpdated(ctx, update.ChatMember)
	}

	return nil
}

func (h *TelegramUpdateHandler) message(ctx context.Context, message *tgbotapi.Message) error {
	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case consts.TgCmdStart:
			if message.Chat != nil {
				return h.roomEnqueue(ctx, cmd, *message.Chat)
			}
		case consts.TgCmdAdd, consts.TgCmdRumors:
			return h.enqueueMessage(ctx, cmd, message)
		case consts.TgCmdRoom, consts.TgCmdFeed:
			if message.Chat != nil {
				if message.Chat.ID == h.Owner || (message.From != nil && message.From.ID == h.Owner) {
					return h.enqueueMessage(ctx, cmd, message)
				} else {
					h.Emitter.Fire(ctx, consts.EventErrorForbidden, message.Chat.ID)
				}
			}
		}
	}
	return nil
}

func (h *TelegramUpdateHandler) enqueueMessage(ctx context.Context, cmd string, message *tgbotapi.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error due to marshal Message")
		return err
	}
	return h.enqueue(ctx, cmd, payload)
}

func (h *TelegramUpdateHandler) chatMemberUpdated(ctx context.Context, member *tgbotapi.ChatMemberUpdated) error {
	if _, ok := commands[member.NewChatMember.Status]; ok {
		return h.roomEnqueue(ctx, member.NewChatMember.Status, member.Chat)
	}
	return nil
}

func (h *TelegramUpdateHandler) roomEnqueue(ctx context.Context, cmd string, chat tgbotapi.Chat) error {
	payload, err := json.Marshal(chat)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error due to marshal MyChatMember.Chat")
		return err
	}
	return h.enqueue(ctx, cmd, payload)
}

func (h *TelegramUpdateHandler) enqueue(ctx context.Context, cmd string, payload []byte) error {
	if _, err := h.Client.EnqueueContext(ctx, asynq.NewTask(commands[cmd], payload)); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("new_task", commands[cmd]).RawJSON("new_payload", payload).Msg("error due to enqueue task")
		return err
	}
	return nil
}
