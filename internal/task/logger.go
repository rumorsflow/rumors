package task

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
	"os"
)

type asynqLogger struct {
	logger *slog.Logger
}

func (l *asynqLogger) Debug(args ...any) {
	l.logger.Debug(fmt.Sprintf(args[0].(string), args[1:]...))
}

func (l *asynqLogger) Info(args ...any) {
	l.logger.Info(fmt.Sprintf(args[0].(string), args[1:]...))
}

func (l *asynqLogger) Warn(args ...any) {
	l.logger.Warn(fmt.Sprintf(args[0].(string), args[1:]...))
}

func (l *asynqLogger) Error(args ...any) {
	l.logger.Error(fmt.Sprintf(args[0].(string), args[1:]...))
}

func (l *asynqLogger) Fatal(args ...any) {
	l.logger.Error(fmt.Sprintf(args[0].(string), args[1:]...))
	os.Exit(1)
}

func level(ctx context.Context, log *slog.Logger) asynq.LogLevel {
	for a, l := range map[asynq.LogLevel]slog.Level{
		asynq.DebugLevel: slog.LevelDebug,
		asynq.InfoLevel:  slog.LevelInfo,
		asynq.WarnLevel:  slog.LevelWarn,
		asynq.ErrorLevel: slog.LevelError,
	} {
		if log.Enabled(ctx, l) {
			return a
		}
	}
	return asynq.InfoLevel
}
