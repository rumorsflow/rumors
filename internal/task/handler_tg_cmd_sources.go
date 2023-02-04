package task

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"golang.org/x/exp/slog"
)

type HandlerTgCmdSources struct {
	logger    *slog.Logger
	publisher *pubsub.Publisher
}

func (h *HandlerTgCmdSources) ProcessTask(ctx context.Context, task *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	feeds := ctx.Value(ctxFeedsKey{}).([]*entity.Feed)

	hosts := make([]string, 0, len(feeds))
	seen := make(map[string]struct{}, len(feeds))
	for _, feed := range feeds {
		if _, ok := seen[feed.Host]; !ok {
			seen[feed.Host] = struct{}{}
			hosts = append(hosts, feed.Host)
		}
	}

	h.publisher.Telegram(ctx, telegram.Message{
		ChatID: message.Chat.ID,
		View:   telegram.ViewSources,
		Data:   hosts,
	})
	return nil
}
