package handlers

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog/log"
)

type FeedItemSaveHandler struct {
	Storage storage.FeedItemStorage
	Client  *asynq.Client
}

func (h *FeedItemSaveHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var feedItem models.FeedItem
	if err := json.Unmarshal(task.Payload(), &feedItem); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	if err := h.Storage.Save(ctx, feedItem); err != nil {
		l.Error().Err(err).Msg("error due to save feed item")
		return nil
	}

	t := asynq.NewTask(consts.TaskFeedItemGroup, task.Payload())
	q := asynq.Queue(consts.QueueFeedItems)
	g := asynq.Group(feedItem.Domain())

	if _, err := h.Client.EnqueueContext(ctx, t, q, g); err != nil {
		l.Error().Err(err).Msg("error due to enqueue feed item")
		return err
	}

	return nil
}
