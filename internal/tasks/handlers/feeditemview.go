package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/rs/zerolog/log"
)

type FeedItemViewHandler struct {
	Storage storage.FeedItemStorage
	Emitter emitter.Emitter
}

func (h *FeedItemViewHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var message tgbotapi.Message
	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	var link *string
	index, size, rest := Pagination(message.CommandArguments())
	if item := slice.Safe(rest, 0); item != "" {
		link = &item
	}

	data, err := h.Storage.Find(ctx, &storage.FilterFeedItems{Link: link}, index, size)
	if err != nil {
		l.Error().Err(err).Msg("error due to find feed items")
		return err
	}

	group := make(map[string][]models.FeedItem)
	for _, item := range data {
		domain := item.Domain()
		if _, ok := group[domain]; ok {
			group[domain] = append(group[domain], item)
		} else {
			group[domain] = []models.FeedItem{item}
		}
	}

	h.Emitter.Fire(ctx, consts.EventFeedItemViewList, message.Chat.ID, group)

	return nil
}
