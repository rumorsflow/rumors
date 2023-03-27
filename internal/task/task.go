package task

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"golang.org/x/exp/slog"
	"os"
)

const (
	OpClientEnqueue errs.Op = "task.client: enqueue"

	OpMetricsRegister errs.Op = "task.metrics: register"
	OpMetricsClose    errs.Op = "task.metrics: close"

	OpServerStart        errs.Op = "task.server: start"
	OpServerProcessTask  errs.Op = "task.server: process task"
	OpServerParseFeed    errs.Op = "task.server: parse feed link"
	OpServerParseSitemap errs.Op = "task.server: parse sitemap link"
	OpServerParseArticle errs.Op = "task.server: parse article link"

	OpSchedulerStart  errs.Op = "task.scheduler: start"
	OpSchedulerSync   errs.Op = "task.scheduler: sync"
	OpSchedulerAdd    errs.Op = "task.scheduler: add"
	OpSchedulerRemove errs.Op = "task.scheduler: remove"

	OpMarshal   errs.Op = "task: marshal payload"
	OpUnmarshal errs.Op = "task: unmarshal payload"
)

const (
	TgCmdStart  = "start"
	TgCmdRumors = "rumors"
	TgCmdSites  = "sites"
	TgCmdSub    = "sub"
	TgCmdOn     = "on"
	TgCmdOff    = "off"
)

const (
	TelegramPrefix    = "telegram:"
	TelegramCmd       = TelegramPrefix + "cmd:"
	TelegramCmdRumors = TelegramCmd + TgCmdRumors
	TelegramCmdSites  = TelegramCmd + TgCmdSites
	TelegramCmdSub    = TelegramCmd + TgCmdSub
	TelegramCmdOn     = TelegramCmd + TgCmdOn
	TelegramCmdOff    = TelegramCmd + TgCmdOff
	TelegramChat      = TelegramPrefix + "chat:"
	TelegramChatNew   = TelegramChat + "new"
	TelegramChatEdit  = TelegramChat + "edit"
)

func level(log *slog.Logger) asynq.LogLevel {
	for a, l := range map[asynq.LogLevel]slog.Level{
		asynq.DebugLevel: slog.LevelDebug,
		asynq.InfoLevel:  slog.LevelInfo,
		asynq.WarnLevel:  slog.LevelWarn,
		asynq.ErrorLevel: slog.LevelError,
	} {
		if log.Enabled(l) {
			return a
		}
	}
	return asynq.InfoLevel
}

func marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, errs.E(OpMarshal, err)
	}
	return data, nil
}

func unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errs.E(OpUnmarshal, err)
	}
	return nil
}

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
	l.logger.Error("asynq error", fmt.Errorf(args[0].(string), args[1:]...))
}

func (l *asynqLogger) Fatal(args ...any) {
	l.logger.Error("asynq error", fmt.Errorf(args[0].(string), args[1:]...))
	os.Exit(1)
}
