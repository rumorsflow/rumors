package task

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"golang.org/x/exp/slog"
)

type HandlerTgCmdOn struct {
	logger    *slog.Logger
	publisher *pubsub.Publisher
	chatRepo  repository.ReadWriteRepository[*entity.Chat]
}

func (h *HandlerTgCmdOn) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	feeds := ctx.Value(ctxFeedsKey{}).([]*entity.Feed)
	chat := ctx.Value(ctxChatKey{}).(*entity.Chat)

	if message.CommandArguments() == "" {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: message.Chat.ID,
			View:   telegram.ViewError,
			Data:   TgErrMsgRequiredSource,
		})
		return nil
	}

	feeds = filterFeedsByHost(feeds, message.CommandArguments())

	if len(feeds) == 0 {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: message.Chat.ID,
			View:   telegram.ViewError,
			Data:   fmt.Sprintf(TgErrMsgNotFoundSource, message.CommandArguments()),
		})
		return nil
	}

	size := len(feeds)
	if chat.Broadcast != nil {
		size += len(*chat.Broadcast)
	}

	ids := make([]uuid.UUID, 0, size)
	seen := make(map[uuid.UUID]struct{}, size)

	for _, feed := range feeds {
		ids = append(ids, feed.ID)
		seen[feed.ID] = struct{}{}
	}

	if chat.Broadcast != nil {
		for _, id := range *chat.Broadcast {
			if _, ok := seen[id]; ok {
				continue
			}
			ids = append(ids, id)
			seen[id] = struct{}{}
		}
	}

	chat.SetBroadcast(ids)

	if err := h.chatRepo.Save(ctx, chat); err != nil {
		h.logger.Error("error due to save chat", err, "id", chat.ID, "telegram_id", chat.TelegramID)

		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: chat.TelegramID,
			View:   telegram.ViewError,
		})
	} else {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: chat.TelegramID,
			View:   telegram.ViewSuccess,
			Data:   TgSuccessMsgSubscribed,
		})
	}

	return nil
}

func filterFeedsByHost(feeds []*entity.Feed, host string) []*entity.Feed {
	return filterFeeds(feeds, func(feed *entity.Feed) bool {
		return feed.Host == host
	})
}

func filterFeeds(feeds []*entity.Feed, condition func(feed *entity.Feed) bool) []*entity.Feed {
	items := make([]*entity.Feed, 0, len(feeds))
	for _, feed := range feeds {
		if condition(feed) {
			items = append(items, feed)
		}
	}
	return items
}
