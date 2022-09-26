package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/daos"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/pkg/litedb/types"
	"github.com/rs/zerolog"
)

type RoomsAddHandler struct {
	Notification notifications.Notification
	Dao          *daos.Dao
	Log          *zerolog.Logger
}

func (h *RoomsAddHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
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
		model.Deleted = false

		if err = h.Dao.Update(ctx, model); err != nil {
			h.Notification.Err(nil, err)
			return nil
		}
	} else {
		model = &models.Room{
			Id:        chat.ID,
			Type:      chat.Type,
			Title:     chat.Title,
			Broadcast: false,
			Deleted:   false,
			Created:   types.NowDateTime(),
		}
		if err = h.Dao.Insert(ctx, model); err != nil {
			h.Notification.Err(nil, err)
			return nil
		}
	}

	h.Notification.Send(nil, model.Info())

	return nil
}
