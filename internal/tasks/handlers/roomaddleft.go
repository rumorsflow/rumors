package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RoomAddLeftHandler struct {
	Storage storage.RoomStorage
	Client  *asynq.Client
}

func (h *RoomAddLeftHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var chat tgbotapi.Chat
	if err := json.Unmarshal(task.Payload(), &chat); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	title := chat.Title
	if title == "" {
		title = chat.FirstName
		if chat.LastName != "" {
			title += " " + chat.LastName
		}
	}

	room, err := h.Storage.FindByChatId(ctx, chat.ID)
	if err == nil {
		room.Title = title
		room.Deleted = task.Type() == consts.TaskRoomLeft
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		room = models.Room{
			Id:        uuid.NewString(),
			ChatId:    chat.ID,
			Type:      chat.Type,
			Title:     title,
			Broadcast: false,
			Deleted:   task.Type() == consts.TaskRoomLeft,
			CreatedAt: time.Now().UTC(),
		}
	} else {
		l.Error().Err(err).Msg("error due to find room by chat id")
		return nil
	}

	payload, _ := json.Marshal(room)
	if _, err = h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskRoomSave, payload)); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("new_task", consts.TaskRoomSave).RawJSON("new_payload", payload).Msg("error due to enqueue task")
		return err
	}

	return nil
}
