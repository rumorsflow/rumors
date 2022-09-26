package background

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/config"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/logger"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

type App struct {
	cfg       config.AsynqConfig
	log       *zerolog.Logger
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	started   bool
	mu        sync.Mutex
}

func NewApp(cfg config.AsynqConfig, zeroLog *zerolog.Logger) *App {
	redisConnOpt := asynq.RedisClientOpt{
		Network:  cfg.Redis.Network,
		Addr:     cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	level := logger.AsyncLogLevel(zeroLog.GetLevel())
	log := logger.NewAsynqLogger(zeroLog)

	return &App{
		cfg:    cfg,
		log:    zeroLog,
		client: asynq.NewClient(redisConnOpt),
		server: asynq.NewServer(redisConnOpt, asynq.Config{
			Concurrency:      cfg.Server.Concurrency,
			LogLevel:         level,
			Logger:           log,
			ShutdownTimeout:  10 * time.Second,
			GroupMaxDelay:    10 * time.Minute,
			GroupGracePeriod: 2 * time.Minute,
			GroupMaxSize:     50,
			GroupAggregator:  asynq.GroupAggregatorFunc(aggregate),
			Queues: map[string]int{
				"default":   1,
				"broadcast": 2,
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

	if a.started {
		return nil
	}

	a.log.Info().Msg("Start asynq server")
	if err := a.server.Start(handler); err != nil {
		return err
	}

	a.log.Info().Msg("Register scheduler")
	task := asynq.NewTask(a.cfg.Scheduler.TaskName, nil)
	entryId, err := a.scheduler.Register(a.cfg.Scheduler.Cronspec, task)
	if err != nil {
		a.server.Shutdown()
		return err
	}
	a.log.Info().Str("entry_id", entryId).Msg("Registered an entry")
	a.log.Info().Msg("Start asynq scheduler")
	if err = a.scheduler.Start(); err != nil {
		a.server.Shutdown()
		return err
	}

	a.started = true

	return nil
}

func (a *App) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.started {
		return
	}

	a.log.Info().Msg("Stop asynq server")
	a.server.Stop()
}

func (a *App) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.started {
		return
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func(scheduler *asynq.Scheduler, log *zerolog.Logger) {
		defer wg.Done()
		a.log.Info().Msg("Shutdown asynq scheduler")
		scheduler.Shutdown()
	}(a.scheduler, a.log)

	go func(server *asynq.Server, log *zerolog.Logger) {
		defer wg.Done()
		a.log.Info().Msg("Shutdown asynq server")
		server.Stop()
		server.Shutdown()
	}(a.server, a.log)

	wg.Wait()

	a.started = false
}

func aggregate(_ string, tasks []*asynq.Task) *asynq.Task {
	var items []models.FeedItem

	for _, task := range tasks {
		var item models.FeedItem
		if err := json.Unmarshal(task.Payload(), &item); err == nil {
			items = append(items, item)
		}
	}

	data, _ := json.Marshal(items)

	return asynq.NewTask("aggregated:broadcast", data)
}
