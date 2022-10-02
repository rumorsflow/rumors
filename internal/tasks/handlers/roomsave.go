package handlers

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/rs/zerolog/log"
)

type RoomSaveHandler struct {
	Storage storage.RoomStorage
	Emitter emitter.Emitter
}

func (h *RoomSaveHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var room models.Room
	if err := json.Unmarshal(task.Payload(), &room); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	if err := h.Storage.Save(ctx, room); err != nil {
		l.Error().Err(err).Msg("error due to save room")
		h.Emitter.Fire(l.WithContext(ctx), consts.EventRoomSaveError, room, err)
		return nil
	}

	h.Emitter.Fire(l.WithContext(ctx), consts.EventRoomSaveAfter, room)

	return nil
}
