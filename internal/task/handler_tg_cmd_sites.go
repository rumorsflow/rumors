package task

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"golang.org/x/exp/slog"
)

type HandlerTgCmdSites struct {
	logger    *slog.Logger
	publisher common.Pub
}

func (h *HandlerTgCmdSites) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	sites := ctx.Value(ctxSitesKey{}).([]*entity.Site)

	domains := make([]string, len(sites))
	for i, site := range sites {
		domains[i] = site.Domain
	}

	h.publisher.Telegram(ctx, model.Message{
		ChatID: message.Chat.ID,
		View:   model.ViewSites,
		Data:   domains,
	})
	return nil
}
