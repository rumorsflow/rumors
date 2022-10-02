package logger

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	golog "log"
	"os"
	"time"
)

func init() {
	log.Logger = log.Logger.Level(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	initAll()
}

func initAll() {
	zerolog.DefaultContextLogger = &log.Logger
	golog.SetOutput(log.Logger)
}

func ZeroLogWith(lvl, cmd, version string, console, debug bool) {
	l := log.Logger

	if debug {
		l = l.Level(zerolog.DebugLevel)
	} else if level, _ := zerolog.ParseLevel(lvl); zerolog.NoLevel != level {
		l = l.Level(level)
	}

	if console {
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Logger = l.With().Str("command", cmd).Str("version", version).Logger()

	initAll()
}

type asynqLogger struct {
	log zerolog.Logger
}

func AsyncLogLevel() asynq.LogLevel {
	switch log.Logger.GetLevel() {
	case zerolog.DebugLevel, zerolog.TraceLevel:
		return asynq.DebugLevel
	case zerolog.InfoLevel:
		return asynq.InfoLevel
	case zerolog.WarnLevel:
		return asynq.WarnLevel
	case zerolog.ErrorLevel:
		return asynq.ErrorLevel
	case zerolog.FatalLevel, zerolog.PanicLevel:
		return asynq.FatalLevel
	}
	return asynq.LogLevel(0)
}

func NewAsynqLogger() asynq.Logger {
	return &asynqLogger{log: log.Logger.With().Str("context", "asynq").Logger()}
}

func (l *asynqLogger) Debug(args ...any) {
	l.log.Debug().Msg(args[0].(string))
}

func (l *asynqLogger) Info(args ...any) {
	l.log.Info().Msg(args[0].(string))
}

func (l *asynqLogger) Warn(args ...any) {
	l.log.Warn().Msg(args[0].(string))
}

func (l *asynqLogger) Error(args ...any) {
	l.log.Error().Msg(args[0].(string))
}

func (l *asynqLogger) Fatal(args ...any) {
	l.log.Fatal().Msg(args[0].(string))
}

type botLogger struct {
	log zerolog.Logger
}

func NewBotLogger() tgbotapi.BotLogger {
	return &botLogger{log: log.Logger.With().Str("context", "tgbotapi").Logger()}
}

func (l *botLogger) Println(v ...any) {
	if len(v) > 0 {
		switch m := v[0].(type) {
		case error:
			l.log.Error().Err(m).Msg("")
		case string:
			l.log.Info().Msg(m)
		}
	}
}

func (l *botLogger) Printf(format string, v ...any) {
	l.log.Info().Msg(fmt.Sprintf(format, v...))
}
