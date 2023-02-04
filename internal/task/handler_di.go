package task

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
)

type ServerMuxKey struct{}

func GetServerMux(ctx context.Context, c ...di.Container) (*asynq.ServeMux, error) {
	return di.Get[*asynq.ServeMux](ctx, ServerMuxKey{}, c...)
}

func ServerMuxActivator() *di.Activator {
	return &di.Activator{
		Key: ServerMuxKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			log := logger.WithGroup("task").WithGroup("server").WithGroup("mux")

			mux := asynq.NewServeMux()
			mux.Use(LoggingMiddleware(log))

			feedRepo, err := db.GetFeedRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			chatRepo, err := db.GetChatRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			articleRepo, err := db.GetArticleRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			publisher, err := pubsub.GetPub(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			hLog := log.WithGroup("handler")

			mux.Handle(JobFeed, &HandlerJobFeed{
				logger:      hLog.WithGroup("job").WithGroup("feed"),
				publisher:   publisher,
				feedRepo:    feedRepo,
				articleRepo: articleRepo,
			})

			tgLog := hLog.WithGroup("telegram")

			mux.Handle(TelegramChat, &HandlerTgChat{
				logger:    tgLog.WithGroup("chat"),
				publisher: publisher,
				chatRepo:  chatRepo,
			})

			cmdLog := tgLog.WithGroup("cmd")

			cmd := asynq.NewServeMux()
			cmd.Use(TgCmdMiddleware(feedRepo, chatRepo, publisher, cmdLog))
			cmd.Handle(TelegramCmdRumors, &HandlerTgCmdRumors{
				logger:      cmdLog.WithGroup("rumors"),
				publisher:   publisher,
				articleRepo: articleRepo,
			})
			cmd.Handle(TelegramCmdSources, &HandlerTgCmdSources{
				logger:    cmdLog.WithGroup("sources"),
				publisher: publisher,
			})
			cmd.Handle(TelegramCmdSub, &HandlerTgCmdSub{
				logger:    cmdLog.WithGroup("sub"),
				publisher: publisher,
			})
			cmd.Handle(TelegramCmdOn, &HandlerTgCmdOn{
				logger:    cmdLog.WithGroup("on"),
				publisher: publisher,
				chatRepo:  chatRepo,
			})
			cmd.Handle(TelegramCmdOff, &HandlerTgCmdOff{
				logger:    cmdLog.WithGroup("off"),
				publisher: publisher,
				chatRepo:  chatRepo,
			})

			mux.Handle(TelegramCmd, cmd)

			return mux, nil, nil
		}),
	}
}
