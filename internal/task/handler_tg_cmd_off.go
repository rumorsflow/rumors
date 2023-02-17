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

type HandlerTgCmdOff struct {
	logger    *slog.Logger
	publisher *pubsub.Publisher
	chatRepo  repository.ReadWriteRepository[*entity.Chat]
}

func (h *HandlerTgCmdOff) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	sites := ctx.Value(ctxSitesKey{}).([]*entity.Site)
	chat := ctx.Value(ctxChatKey{}).(*entity.Chat)

	if message.CommandArguments() == "" {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: message.Chat.ID,
			View:   telegram.ViewError,
			Data:   TgErrMsgRequiredSite,
		})
		return nil
	}

	sites = filterSitesByDomain(sites, message.CommandArguments())

	if len(sites) == 0 {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: message.Chat.ID,
			View:   telegram.ViewError,
			Data:   fmt.Sprintf(TgErrMsgNotFoundSite, message.CommandArguments()),
		})
		return nil
	}

	if chat.Broadcast == nil || len(*chat.Broadcast) == 0 {
		h.publisher.Telegram(ctx, telegram.Message{
			ChatID: chat.TelegramID,
			View:   telegram.ViewError,
		})
		return nil
	}

	ids := make([]uuid.UUID, 0, len(*chat.Broadcast))
	seen := make(map[uuid.UUID]struct{}, len(sites))

	for _, site := range sites {
		seen[site.ID] = struct{}{}
	}

	for _, id := range *chat.Broadcast {
		if _, ok := seen[id]; ok {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) != len(*chat.Broadcast) {
		chat.SetBroadcast(ids)
		if err := h.chatRepo.Save(ctx, chat); err != nil {
			h.logger.Error("error due to save chat", err, "id", chat.ID, "telegram_id", chat.TelegramID)

			h.publisher.Telegram(ctx, telegram.Message{
				ChatID: chat.TelegramID,
				View:   telegram.ViewError,
			})
			return nil
		}
	}

	h.publisher.Telegram(ctx, telegram.Message{
		ChatID: chat.TelegramID,
		View:   telegram.ViewSuccess,
		Data:   TgSuccessMsgUnsubscribed,
	})

	return nil
}
