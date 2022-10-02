package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/rs/zerolog/log"
)

type FeedViewHandler struct {
	Storage storage.FeedStorage
	Emitter emitter.Emitter
}

func (h *FeedViewHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message
	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	if _, err := uuid.Parse(message.CommandArguments()); err == nil {
		return h.one(ctx, message)
	}
	return h.list(ctx, message)
}

func (h *FeedViewHandler) one(ctx context.Context, message tgbotapi.Message) error {
	feed, err := h.Storage.FindById(ctx, message.CommandArguments())
	if err == nil {
		h.Emitter.Fire(ctx, consts.EventFeedViewOne, message.Chat.ID, feed)
	} else {
		log.Ctx(ctx).Error().Err(err).Msg("error due to find feed by id")
		h.Emitter.Fire(ctx, consts.EventErrorNotFound, message.Chat.ID, err)
	}
	return nil
}

func (h *FeedViewHandler) list(ctx context.Context, message tgbotapi.Message) error {
	var link *string
	index, size, rest := Pagination(message.CommandArguments())
	if item := slice.Safe(rest, 0); item != "" {
		link = &item
	}

	feeds, err := h.Storage.Find(ctx, &storage.FilterFeeds{Link: link}, index, size)
	if err == nil {
		h.Emitter.Fire(ctx, consts.EventFeedViewList, message.Chat.ID, feeds)
	} else {
		log.Ctx(ctx).Error().Err(err).Msg("error due to find feed list")
		h.Emitter.Fire(ctx, consts.EventErrorViewList, message.Chat.ID, err)
	}
	return nil
}
