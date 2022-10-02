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

type FeedSaveHandler struct {
	Storage storage.FeedStorage
	Emitter emitter.Emitter
}

func (h *FeedSaveHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var feed models.Feed
	if err := json.Unmarshal(task.Payload(), &feed); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	if err := h.Storage.Save(ctx, feed); err != nil {
		l.Error().Err(err).Msg("error due to save feed")
		h.Emitter.Fire(l.WithContext(ctx), consts.EventFeedSaveError, feed, err)
		return nil
	}

	h.Emitter.Fire(l.WithContext(ctx), consts.EventFeedSaveAfter, feed)

	return nil
}
