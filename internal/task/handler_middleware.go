package task

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
)

const (
	TgSuccessMsgSubscribed   = "Subscribed successfully."
	TgSuccessMsgUnsubscribed = "Unsubscribed successfully."

	TgErrMsgRequiredSite = "Site (domain) is required."
	TgErrMsgNotFoundSite = "Site `%s` not found."
)

type (
	ctxMsgKey   struct{}
	ctxChatKey  struct{}
	ctxSitesKey struct{}
)

func LoggingMiddleware(log *slog.Logger) asynq.MiddlewareFunc {
	return func(handler asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			if logger.IsDebug() {
				log.Debug("handle task", "task", task.Type(), "payload", task.Payload())
			}

			return handler.ProcessTask(ctx, task)
		})
	}
}

func TgCmdMiddleware(
	siteRepo repository.ReadRepository[*entity.Site],
	chatRepo repository.ReadWriteRepository[*entity.Chat],
	publisher *pubsub.Publisher,
	logger *slog.Logger,
) asynq.MiddlewareFunc {
	return func(handler asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			var message tgbotapi.Message
			if err := unmarshal(task.Payload(), &message); err != nil {
				logger.Error("error due to unmarshal task payload", err)
				return nil
			}

			logger.Info("task processing", "telegram_id", message.Chat.ID, "command", message.Command(), "args", message.CommandArguments())

			criteria := db.BuildCriteria(fmt.Sprintf("size=1&field.0.0=telegram_id&value.0.0=%d", message.Chat.ID))
			chats, err := chatRepo.Find(ctx, criteria)
			if err != nil {
				logger.Error("error due to find chat", err, "command", message.Command(), "telegram_id", message.Chat.ID)
				return err
			}
			if len(chats) == 0 {
				logger.Warn("error due to chat not found", "command", message.Command(), "telegram_id", message.Chat.ID)
				return nil
			}
			if chats[0].IsBlocked() {
				logger.Warn("error due to chat was blocked", "command", message.Command(), "telegram_id", message.Chat.ID)
				return nil
			}

			sites, err := siteRepo.Find(ctx, db.BuildCriteria("sort=domain&field.0.0=enabled&value.0.0=true"))
			if err != nil {
				err = fmt.Errorf("%s error: %w", OpServerProcessTask, err)
				logger.Error("error due to find sites", err, "command", message.Command())
				return err
			}

			if len(sites) == 0 {
				logger.Warn("task processing was stopped, because no sites found", "command", message.Command(), "args", message.CommandArguments())
				publisher.Telegram(ctx, telegram.Message{ChatID: message.Chat.ID, View: telegram.ViewError})
				return nil
			}

			ctx = context.WithValue(ctx, ctxMsgKey{}, message)
			ctx = context.WithValue(ctx, ctxChatKey{}, chats[0])
			ctx = context.WithValue(ctx, ctxSitesKey{}, sites)

			return handler.ProcessTask(ctx, task)
		})
	}
}
