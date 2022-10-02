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
	mu        sync.Mutex
}

func NewApp(cfg config.AsynqConfig) *App {
	redisConnOpt := asynq.RedisClientOpt{
		Network:  cfg.Redis.Network,
		Addr:     cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

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
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.client
}

func (a *App) Start(handler asynq.Handler) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.log.Info("Start asynq server")
	if err := a.server.Start(handler); err != nil {
		return err
	}

	a.log.Info("Register scheduler")
	task := asynq.NewTask(consts.TaskFeedScheduler, nil)
	entryId, err := a.scheduler.Register(a.cfg.Scheduler.FeedImporter, task)
	if err != nil {
		a.server.Shutdown()
		return err
	}
	a.log.Info(fmt.Sprintf("scheduler entry %s was registered", entryId))
	a.log.Info("Start asynq scheduler")
	if err = a.scheduler.Start(); err != nil {
		a.server.Shutdown()
		return err
	}

	return nil
}

func (a *App) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

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
