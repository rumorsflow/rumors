package task

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"golang.org/x/exp/slog"
)

type HandlerTgCmdSub struct {
	logger    *slog.Logger
	publisher common.Pub
}

func (h *HandlerTgCmdSub) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	sites := ctx.Value(ctxSitesKey{}).([]*entity.Site)
	chat := ctx.Value(ctxChatKey{}).(*entity.Chat)

	if chat.Broadcast == nil || len(*chat.Broadcast) == 0 {
		h.publisher.Telegram(ctx, model.Message{
			ChatID: message.Chat.ID,
			View:   model.ViewNotFound,
		})
		return nil
	}

	b := make(map[uuid.UUID]struct{}, len(*chat.Broadcast))
	for _, id := range *chat.Broadcast {
		b[id] = struct{}{}
	}

	domains := make([]string, 0, len(b))
	for _, site := range sites {
		if _, ok := b[site.ID]; ok {
			domains = append(domains, site.Domain)
		}
	}

	h.publisher.Telegram(ctx, model.Message{
		ChatID: message.Chat.ID,
		View:   model.ViewSub,
		Data:   domains,
	})

	return nil
}
