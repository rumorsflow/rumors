package tasks

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/config"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/pkg/logger"
	"sync"
	"time"
)

type App struct {
	cfg       config.AsynqConfig
	log       asynq.Logger
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
}

func NewApp(cfg config.AsynqConfig, redisConnOpt asynq.RedisClientOpt) *App {
	level := logger.AsyncLogLevel()
	log := logger.NewAsynqLogger()

	return &App{
		cfg:    cfg,
		log:    log,
		client: asynq.NewClient(redisConnOpt),
		server: asynq.NewServer(redisConnOpt, asynq.Config{
			ShutdownTimeout:  10 * time.Second,
			Concurrency:      cfg.Server.Concurrency,
			LogLevel:         level,
			Logger:           log,
			GroupMaxSize:     cfg.Server.Group.Max.Size,
			GroupMaxDelay:    cfg.Server.Group.Max.Delay,
			GroupGracePeriod: cfg.Server.Group.Grace.Period,
			GroupAggregator:  &Aggregator{Log: log},
			Queues: map[string]int{
				consts.QueueDefault:   1,
				consts.QueueFeedItems: 2,
			},
		}),
		scheduler: asynq.NewScheduler(redisConnOpt, &asynq.SchedulerOpts{
			LogLevel: level,
			Logger:   log,
		}),
	}
}

func (a *App) Client() *asynq.Client {
	return a.client
}

func (a *App) Start(handler asynq.Handler) {
	a.log.Info("Start asynq server")
	if err := a.server.Start(handler); err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Register scheduler")
	task := asynq.NewTask(consts.TaskFeedScheduler, nil)
	entryId, err := a.scheduler.Register(a.cfg.Scheduler.FeedImporter, task)
	if err != nil {
		a.log.Error(err.Error())
		a.server.Shutdown()
	}
	a.log.Info(fmt.Sprintf("scheduler entry %s was registered", entryId))
	a.log.Info("Start asynq scheduler")
	if err = a.scheduler.Start(); err != nil {
		a.log.Error(err.Error())
		a.server.Shutdown()
	}
}

func (a *App) Shutdown() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(scheduler *asynq.Scheduler, log asynq.Logger) {
		defer wg.Done()
		log.Info("Shutdown asynq scheduler")
		scheduler.Shutdown()
	}(a.scheduler, a.log)

	go func(server *asynq.Server, log asynq.Logger) {
		defer wg.Done()
		log.Info("Shutdown asynq server")
		server.Stop()
		server.Shutdown()
	}(a.server, a.log)

	wg.Wait()
}
