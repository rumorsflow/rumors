package handlers

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/rs/zerolog/log"
)

type FeedItemBroadcastHandler struct {
	Emitter emitter.Emitter
}

func (h *FeedItemBroadcastHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var items []models.FeedItem
	if err := json.Unmarshal(task.Payload()[8:], &items); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	chatId := BytesToInt64(task.Payload()[:8])
	group := map[string][]models.FeedItem{
		items[0].Domain(): items,
	}

	h.Emitter.Fire(l.WithContext(ctx), consts.EventFeedItemViewList, chatId, group)
	return nil
}
