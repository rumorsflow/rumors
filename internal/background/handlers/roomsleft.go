package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/daos"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/rs/zerolog"
)

type RoomsLeftHandler struct {
	Notification notifications.Notification
	Dao          *daos.Dao
	Log          *zerolog.Logger
}

func (h *RoomsLeftHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var chat tgbotapi.Chat

	if err := json.Unmarshal(task.Payload(), &chat); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	model, err := h.Dao.FindRoomById(ctx, chat.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	if model != nil {
		model.Title = chat.Title
		model.Deleted = true

		if err = h.Dao.Update(ctx, model); err != nil {
			h.Notification.Err(nil, err)
			return nil
		}

		h.Notification.Send(nil, model.Info())
	}

	return nil
}
