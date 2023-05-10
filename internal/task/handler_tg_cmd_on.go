package task

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"golang.org/x/exp/slog"
	"strings"
)

type HandlerTgCmdOn struct {
	logger    *slog.Logger
	publisher common.Pub
	chatRepo  repository.ReadWriteRepository[*entity.Chat]
}

func (h *HandlerTgCmdOn) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	sites := ctx.Value(ctxSitesKey{}).([]*entity.Site)
	chat := ctx.Value(ctxChatKey{}).(*entity.Chat)

	if message.CommandArguments() == "" {
		h.publisher.Telegram(ctx, model.Message{
			ChatID: message.Chat.ID,
			View:   model.ViewError,
			Data:   TgErrMsgRequiredSite,
		})
		return nil
	}

	sites = filterSitesByDomain(sites, message.CommandArguments())

	if len(sites) == 0 {
		h.publisher.Telegram(ctx, model.Message{
			ChatID: message.Chat.ID,
			View:   model.ViewError,
			Data:   fmt.Sprintf(TgErrMsgNotFoundSite, message.CommandArguments()),
		})
		return nil
	}

	size := len(sites)
	if chat.Broadcast != nil {
		size += len(*chat.Broadcast)
	}

	ids := make([]uuid.UUID, 0, size)
	seen := make(map[uuid.UUID]struct{}, size)

	for _, site := range sites {
		ids = append(ids, site.ID)
		seen[site.ID] = struct{}{}
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
		h.logger.Error("error due to save chat", "err", err, "id", chat.ID, "telegram_id", chat.TelegramID)

		h.publisher.Telegram(ctx, model.Message{
			ChatID: chat.TelegramID,
			View:   model.ViewError,
		})
	} else {
		h.publisher.Telegram(ctx, model.Message{
			ChatID: chat.TelegramID,
			View:   model.ViewSuccess,
			Data:   TgSuccessMsgSubscribed,
		})
	}

	return nil
}

func filterSitesByDomain(sites []*entity.Site, domain string) []*entity.Site {
	return filterSites(sites, func(site *entity.Site) bool {
		return strings.Contains(site.Domain, domain)
	})
}

func filterSites(sites []*entity.Site, condition func(site *entity.Site) bool) []*entity.Site {
	items := make([]*entity.Site, 0, len(sites))
	for _, site := range sites {
		if condition(site) {
			items = append(items, site)
		}
	}
	return items
}
