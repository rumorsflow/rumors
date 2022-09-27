package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RoomsAddHandler struct {
	Notification notifications.Notification
	RoomStorage  storage.RoomStorage
	Log          *zerolog.Logger
}

func (h *RoomsAddHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var chat tgbotapi.Chat

	if err := json.Unmarshal(task.Payload(), &chat); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	model, err := h.RoomStorage.FindByChatId(ctx, chat.ID)
	if err == nil {
		model.Title = chat.Title
		model.Deleted = false
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		h.Log.Error().Err(err).Msg("")
		return nil
	} else {
		model = models.Room{
			Id:        uuid.NewString(),
			ChatId:    chat.ID,
			Type:      chat.Type,
			Title:     chat.Title,
			Broadcast: false,
			Deleted:   false,
			CreatedAt: time.Now().UTC(),
		}
	}

	if err = h.RoomStorage.Save(ctx, model); err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	h.Notification.Send(nil, model.Info())

	return nil
}
