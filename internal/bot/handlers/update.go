package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/rs/zerolog"
)

type UpdateHandler struct {
	Notification notifications.Notification
	Client       *asynq.Client
	Log          *zerolog.Logger
	Owner        int64
}

var routes = map[string]string{
	"start":         "rooms:add",
	"room":          "rooms:crud",
	"creator":       "rooms:add",
	"member":        "rooms:add",
	"administrator": "rooms:add",
	"restricted":    "rooms:left",
	"left":          "rooms:left",
	"kicked":        "rooms:left",
	"add":           "feeds:add",
	"feed":          "feeds:crud",
	"rumors":        "rumors:list",
}

func (h *UpdateHandler) Process(update tgbotapi.Update) error {
	if update.Message != nil {
		return h.message(update.UpdateID, update.Message)
	}

	if update.EditedMessage != nil {
		return h.message(update.UpdateID, update.EditedMessage)
	}

	if update.ChannelPost != nil {
		return h.message(update.UpdateID, update.ChannelPost)
	}

	if update.EditedChannelPost != nil {
		return h.message(update.UpdateID, update.EditedChannelPost)
	}

	if update.MyChatMember != nil {
		return h.myChatMember(update.UpdateID, update.MyChatMember)
	}

	data, _ := json.Marshal(update)
	h.Log.Info().RawJSON("tg_update", data).Msg("")

	return nil
}

func (h *UpdateHandler) message(updateId int, message *tgbotapi.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			if message.Chat != nil {
				return h.roomEnqueue(updateId, cmd, *message.Chat)
			}
		case "add", "rumors":
			return h.enqueue(updateId, cmd, data)
		case "room", "feed":
			if message.Chat.ID == h.Owner || (message.From != nil && message.From.ID == h.Owner) {
				return h.enqueue(updateId, cmd, data)
			}
			h.Notification.Forbidden(message.Chat.ID)
			return nil
		}
	}

	h.Log.Info().RawJSON("tg_update_message", data).Msg("")

	return nil
}

func (h *UpdateHandler) myChatMember(updateId int, member *tgbotapi.ChatMemberUpdated) error {
	if _, ok := routes[member.NewChatMember.Status]; ok {
		return h.roomEnqueue(updateId, member.NewChatMember.Status, member.Chat)
	}
	return nil
}

func (h *UpdateHandler) roomEnqueue(updateId int, name string, chat tgbotapi.Chat) error {
	data, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	h.Log.Info().RawJSON("tg_update_room", data).Msg("")

	return h.enqueue(updateId, name, data)
}

func (h *UpdateHandler) enqueue(updateId int, cmd string, data []byte) error {
	task := asynq.NewTask(routes[cmd], data)
	taskId := asynq.TaskID(fmt.Sprintf("%s:%d", routes[cmd], updateId))

	_, err := h.Client.Enqueue(task, taskId)
	if errors.Is(err, asynq.ErrTaskIDConflict) {
		h.Log.Warn().Err(err).Msg("")
		return nil
	}
	return err
}
