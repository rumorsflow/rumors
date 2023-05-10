package task

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

type Server struct {
	cfg *asynq.Config
	srv *asynq.Server
}

type ServerOption func(s *Server)

func WithGroupAggregator(groupAggregator asynq.GroupAggregator) ServerOption {
	return func(s *Server) {
		s.cfg.GroupAggregator = groupAggregator
	}
}

func WithErrorHandler(errorHandler asynq.ErrorHandler) ServerOption {
	return func(s *Server) {
		s.cfg.ErrorHandler = errorHandler
	}
}

func NewServer(cfg *ServerConfig, redisConnOpt asynq.RedisConnOpt, logger *slog.Logger, options ...ServerOption) *Server {
	srv := &Server{
		cfg: &asynq.Config{
			Concurrency:              cfg.Concurrency,
			Queues:                   cfg.Queues,
			StrictPriority:           cfg.StrictPriority,
			HealthCheckInterval:      cfg.HealthCheckInterval,
			DelayedTaskCheckInterval: cfg.DelayedTaskCheckInterval,
			GroupGracePeriod:         cfg.GroupGracePeriod,
			GroupMaxDelay:            cfg.GroupMaxDelay,
			GroupMaxSize:             cfg.GroupMaxSize,
			ShutdownTimeout:          cfg.GracefulTimeout,
			Logger:                   &asynqLogger{logger: logger},
			LogLevel:                 level(context.Background(), logger),
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("handle task error", "err", err, "task", task.Type(), "payload", task.Payload())
			}),
		},
	}

	for _, option := range options {
		option(srv)
	}

	srv.srv = asynq.NewServer(redisConnOpt, *srv.cfg)

	return srv
}

func (s *Server) Start(handler asynq.Handler, errCh chan<- error) {
	if err := s.srv.Start(handler); err != nil {
		errCh <- fmt.Errorf("%s %w", OpServerStart, err)
	}
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.srv.Shutdown()
}
