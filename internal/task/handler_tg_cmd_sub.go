package task

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"golang.org/x/exp/slog"
)

type HandlerTgCmdSub struct {
	logger    *slog.Logger
	publisher *pubsub.Publisher
}

func (h *HandlerTgCmdSub) ProcessTask(ctx context.Context, task *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	feeds := ctx.Value(ctxFeedsKey{}).([]*entity.Feed)
	chat := ctx.Value(ctxChatKey{}).(*entity.Chat)

	if chat.Broadcast == nil || len(*chat.Broadcast) == 0 {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: message.Chat.ID,
			View:   telegram.ViewNotFound,
		})
		return nil
	}

	b := make(map[uuid.UUID]struct{}, len(*chat.Broadcast))
	for _, id := range *chat.Broadcast {
		b[id] = struct{}{}
	}

	hosts := make([]string, 0, len(*chat.Broadcast))
	seen := make(map[string]struct{}, len(*chat.Broadcast))
	for _, feed := range feeds {
		if _, ok := b[feed.ID]; ok {
			if _, ok := seen[feed.Host]; !ok {
				seen[feed.Host] = struct{}{}
				hosts = append(hosts, feed.Host)
			}
		}
	}

	h.publisher.Telegram(ctx, telegram.Message{
		ChatID: message.Chat.ID,
		View:   telegram.ViewSub,
		Data:   hosts,
	})

	return nil
}
