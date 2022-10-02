package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/iagapie/rumors/pkg/validate"
	"github.com/rs/zerolog/log"
)

type FeedAddHandler struct {
	Validator validate.Validator
	Emitter   emitter.Emitter
	Client    *asynq.Client
	Owner     int64
}

func (h *FeedAddHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	l := log.Ctx(ctx)

	var message tgbotapi.Message
	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	args := Args(message.CommandArguments())
	feed := models.Feed{
		By:   message.Chat.ID,
		Link: slice.Safe(args, 0),
		Lang: slice.Safe(args, 1),
	}
	feed.Enabled = feed.By == h.Owner && feed.Lang != ""

	if err := h.Validator.Validate(feed); err != nil {
		h.Emitter.Fire(ctx, consts.EventErrorArgs, message.Chat.ID, "Feed URL is required.")
		return nil
	}

	payload, _ := json.Marshal(feed)

	if _, err := h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskFeedImporter, payload)); err != nil {
		l.Error().Err(err).RawJSON("feed", payload).Msg("error due to enqueue feed")
		return err
	}

	return nil
}
