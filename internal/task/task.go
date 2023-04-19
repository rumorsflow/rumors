package task

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
	"os"
)

const (
	OpClientEnqueue = "task.client: enqueue ->"

	OpMetricsRegister = "task.metrics: register ->"
	OpMetricsClose    = "task.metrics: close ->"

	OpServerStart        = "task.server: start ->"
	OpServerProcessTask  = "task.server: process task ->"
	OpServerParseFeed    = "task.server: parse feed link ->"
	OpServerParseSitemap = "task.server: parse sitemap link ->"
	OpServerParseArticle = "task.server: parse article link ->"

	OpSchedulerStart  = "task.scheduler: start ->"
	OpSchedulerSync   = "task.scheduler: sync ->"
	OpSchedulerAdd    = "task.scheduler: add ->"
	OpSchedulerRemove = "task.scheduler: remove ->"

	OpMarshal   = "task: marshal payload ->"
	OpUnmarshal = "task: unmarshal payload ->"
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
		return nil, fmt.Errorf("%s error: %w", OpMarshal, err)
	}
	return data, nil
}

func unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s error: %w", OpUnmarshal, err)
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
