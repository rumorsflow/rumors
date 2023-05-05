package task

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"golang.org/x/exp/slog"
)

type Server struct {
	cfg             *ServerConfig
	logger          *slog.Logger
	rdbMaker        *rdb.UniversalClientMaker
	groupAggregator asynq.GroupAggregator
	errorHandler    asynq.ErrorHandler
}

type ServerOption func(s *Server)

func WithGroupAggregator(groupAggregator asynq.GroupAggregator) ServerOption {
	return func(s *Server) {
		s.groupAggregator = groupAggregator
	}
}

func WithErrorHandler(errorHandler asynq.ErrorHandler) ServerOption {
	return func(s *Server) {
		s.errorHandler = errorHandler
	}
}

func NewServer(cfg *ServerConfig, rdbMaker *rdb.UniversalClientMaker, options ...ServerOption) *Server {
	cfg.Init()

	srv := &Server{cfg: cfg, rdbMaker: rdbMaker, logger: logger.WithGroup("task").WithGroup("server")}

	for _, option := range options {
		option(srv)
	}

	return srv
}

func (s *Server) Run(ctx context.Context, handler asynq.Handler) error {
	cfg := asynq.Config{
		Concurrency:              s.cfg.Concurrency,
		Queues:                   s.cfg.Queues,
		StrictPriority:           s.cfg.StrictPriority,
		HealthCheckInterval:      s.cfg.HealthCheckInterval,
		DelayedTaskCheckInterval: s.cfg.DelayedTaskCheckInterval,
		GroupGracePeriod:         s.cfg.GroupGracePeriod,
		GroupMaxDelay:            s.cfg.GroupMaxDelay,
		GroupMaxSize:             s.cfg.GroupMaxSize,
		ShutdownTimeout:          s.cfg.GracefulTimeout,
		GroupAggregator:          s.groupAggregator,
		ErrorHandler:             s.errorHandler,
		BaseContext: func() context.Context {
			return ctx
		},
	}

	if s.logger != nil {
		cfg.Logger = &asynqLogger{logger: s.logger}
		cfg.LogLevel = level(s.logger)

		if s.errorHandler == nil {
			s.errorHandler = asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				s.logger.Error("handle task error", "err", err, "task", task.Type(), "payload", task.Payload())
			})
		}
	}

	srv := asynq.NewServer(s.rdbMaker, cfg)

	if err := srv.Start(handler); err != nil {
		return fmt.Errorf("%s error: %w", OpServerStart, err)
	}

	defer func() {
		srv.Stop()
		srv.Shutdown()
	}()

	<-ctx.Done()

	return nil
}
