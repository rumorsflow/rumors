package task

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"golang.org/x/sync/errgroup"
)

const (
	PluginName = "task"

	sectionScheduler = "task.scheduler"
	sectionServer    = "task.server"
)

type Plugin struct {
	client    *Client
	server    *Server
	scheduler *Scheduler
	metrics   *Metrics
	handler   asynq.Handler
}

func (p *Plugin) Init(
	cfg config.Configurer,
	uow common.UnitOfWork,
	redisConnOpt asynq.RedisConnOpt,
	pub common.Pub,
	log logger.Logger,
) error {
	const op = errors.Op("task_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	l := log.NamedLogger(PluginName)

	p.client = NewClient(redisConnOpt, l.WithGroup("client"))

	if cfg.Has(sectionServer) {
		var c ServerConfig
		if err := cfg.UnmarshalKey(sectionServer, &c); err != nil {
			return errors.E(op, err)
		}
		c.Init()
		if c.GracefulTimeout == 0 {
			c.GracefulTimeout = cfg.GracefulTimeout()
		}

		siteAny, err := uow.Repository((*entity.Site)(nil))
		if err != nil {
			return errors.E(op, err)
		}

		chatAny, err := uow.Repository((*entity.Chat)(nil))
		if err != nil {
			return errors.E(op, err)
		}

		articleAny, err := uow.Repository((*entity.Article)(nil))
		if err != nil {
			return errors.E(op, err)
		}

		siteRepo := siteAny.(repository.ReadWriteRepository[*entity.Site])
		chatRepo := chatAny.(repository.ReadWriteRepository[*entity.Chat])
		articleRepo := articleAny.(repository.ReadWriteRepository[*entity.Article])

		ls := l.WithGroup("server")
		muxLog := ls.WithGroup("mux")
		hLog := muxLog.WithGroup("handler")
		tgLog := hLog.WithGroup("telegram")
		cmdLog := tgLog.WithGroup("cmd")

		p.server = NewServer(&c, redisConnOpt, ls)

		mux := asynq.NewServeMux()
		mux.Use(LoggingMiddleware(muxLog))

		mux.Handle(string(entity.JobFeed), &HandlerJobFeed{
			logger:      hLog.WithGroup("job").WithGroup("feed"),
			publisher:   pub,
			siteRepo:    siteRepo,
			articleRepo: articleRepo,
		})

		mux.Handle(string(entity.JobSitemap), &HandlerJobSitemap{
			logger:      hLog.WithGroup("job").WithGroup("sitemap"),
			publisher:   pub,
			siteRepo:    siteRepo,
			articleRepo: articleRepo,
		})

		mux.Handle(TelegramChat, &HandlerTgChat{
			logger:    tgLog.WithGroup("chat"),
			publisher: pub,
			chatRepo:  chatRepo,
		})

		cmd := asynq.NewServeMux()
		cmd.Use(TgCmdMiddleware(siteRepo, chatRepo, pub, cmdLog))
		cmd.Handle(TelegramCmdRumors, &HandlerTgCmdRumors{
			logger:      cmdLog.WithGroup("rumors"),
			publisher:   pub,
			articleRepo: articleRepo,
		})
		cmd.Handle(TelegramCmdSites, &HandlerTgCmdSites{
			logger:    cmdLog.WithGroup("sites"),
			publisher: pub,
		})
		cmd.Handle(TelegramCmdSub, &HandlerTgCmdSub{
			logger:    cmdLog.WithGroup("sub"),
			publisher: pub,
		})
		cmd.Handle(TelegramCmdOn, &HandlerTgCmdOn{
			logger:    cmdLog.WithGroup("on"),
			publisher: pub,
			chatRepo:  chatRepo,
		})
		cmd.Handle(TelegramCmdOff, &HandlerTgCmdOff{
			logger:    cmdLog.WithGroup("off"),
			publisher: pub,
			chatRepo:  chatRepo,
		})

		mux.Handle(TelegramCmd, cmd)

		p.handler = mux
	}

	if cfg.Has(sectionScheduler) {
		var c SchedulerConfig
		if err := cfg.UnmarshalKey(sectionScheduler, &c); err != nil {
			return errors.E(op, err)
		}
		c.Init()

		jobRepo, err := uow.Repository((*entity.Job)(nil))
		if err != nil {
			return errors.E(op, err)
		}

		p.scheduler = NewScheduler(
			jobRepo.(repository.ReadWriteRepository[*entity.Job]),
			redisConnOpt,
			l.WithGroup("scheduler"),
			WithInterval(c.SyncInterval),
		)
	}

	if p.server != nil || p.scheduler != nil {
		p.metrics = NewMetrics(redisConnOpt, l.WithGroup("metrics"))

		if err := p.metrics.Register(); err != nil {
			return errors.E(op, err)
		}
	}

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 2)

	if p.server != nil {
		p.server.Start(p.handler, errCh)
	}

	if p.scheduler != nil {
		p.scheduler.Start(context.Background(), errCh)
	}

	return errCh
}

func (p *Plugin) Stop(ctx context.Context) (err error) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(p.client.Close)

	if p.metrics != nil {
		g.Go(func() error {
			p.metrics.Unregister()
			return p.metrics.Close()
		})
	}

	if p.scheduler != nil {
		g.Go(func() error {
			p.scheduler.Stop()
			return nil
		})
	}

	if p.server != nil {
		g.Go(func() error {
			p.server.Stop()
			return nil
		})
	}

	return g.Wait()
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Client() common.Client {
	return p.client
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*common.Client)(nil), p.Client),
	}
}
