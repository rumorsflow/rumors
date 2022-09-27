package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomsLeftHandler struct {
	Notification notifications.Notification
	RoomStorage  storage.RoomStorage
	Log          *zerolog.Logger
}

func (h *RoomsLeftHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var chat tgbotapi.Chat

	if err := json.Unmarshal(task.Payload(), &chat); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	model, err := h.RoomStorage.FindByChatId(ctx, chat.ID)
	if err == nil {
		model.Title = chat.Title
		model.Deleted = true

		if err = h.RoomStorage.Save(ctx, model); err != nil {
			h.Notification.Err(nil, err)
			return nil
		}

		h.Notification.Send(nil, model.Info())
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		h.Log.Error().Err(err).Msg("")
	}

	return nil
}
