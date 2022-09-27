package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog"
	"strings"
)

type RoomsHandler struct {
	Notification notifications.Notification
	RoomStorage  storage.RoomStorage
	Client       *asynq.Client
	Log          *zerolog.Logger
}

func (h *RoomsHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message

	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	switch task.Type() {
	case "rooms:crud":
		var name string

		switch strings.ToLower(Args(message.CommandArguments())[0]) {
		case "view", "show", "v", "s":
			name = "rooms:view"
		case "update", "edit", "u", "e", "b":
			name = "rooms:update"
		default:
			name = "rooms:list"
		}

		_, err := h.Client.Enqueue(asynq.NewTask(name, task.Payload()))
		return err
	case "rooms:list":
		return h.list(ctx, message)
	case "rooms:view":
		return h.view(ctx, message)
	case "rooms:update":
		return h.update(ctx, message)
	}

	return nil
}

func (h *RoomsHandler) list(ctx context.Context, message tgbotapi.Message) error {
	i, s, f := Pagination(message.CommandArguments())
	var t *string
	if len(f) > 0 {
		t = &f[0]
	}

	data, err := h.RoomStorage.Find(ctx, storage.FilterRooms{Title: t}, i, s)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	if len(data) == 0 {
		h.Notification.Error(nil, "<b>No Rooms...</b>")
		return nil
	}

	var b strings.Builder
	for j, item := range data {
		b.WriteString(item.Line())

		if (j + 1) < len(data) {
			b.WriteString("\n")
		}
	}

	h.Notification.Send(nil, b.String())
	return nil
}

func (h *RoomsHandler) view(ctx context.Context, message tgbotapi.Message) error {
	id, _ := Id(message.CommandArguments())
	if id == 0 {
		h.Notification.Error(nil, "ID is required")
		return nil
	}

	room, err := h.RoomStorage.FindByChatId(ctx, id)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	h.Notification.Send(nil, room.Info())
	return nil
}

func (h *RoomsHandler) update(ctx context.Context, message tgbotapi.Message) error {
	id, _ := Id(message.CommandArguments())
	if id == 0 {
		h.Notification.Error(nil, "ID is required")
		return nil
	}

	room, err := h.RoomStorage.FindByChatId(ctx, id)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	room.Broadcast = !room.Broadcast
	if err = h.RoomStorage.Save(ctx, room); err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	h.Notification.Send(nil, room.Info())
	return nil
}
