package logger

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"io"
	"os"
	"sync"
	"time"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"
}

type ZeroLog struct {
	Command  string
	Version  string
	LogLevel string
	Colored  bool
	Debug    bool
	once     sync.Once
	log      *zerolog.Logger
}

func (zl *ZeroLog) Log() *zerolog.Logger {
	zl.once.Do(func() {
		level := zerolog.DebugLevel
		if !zl.Debug {
			level, _ = zerolog.ParseLevel(zl.LogLevel)
			if zerolog.NoLevel == level {
				level = zerolog.InfoLevel
			}
		}

		var out io.Writer = os.Stdout
		if zl.Colored {
			out = zerolog.ConsoleWriter{Out: out}
		}

		log := zerolog.New(out).Level(level).With().Timestamp().Str("command", zl.Command).Str("version", zl.Version).Logger()

		zl.log = &log
	})
	return zl.log
}

var _ asynq.Logger = (*asynqLogger)(nil)

type asynqLogger struct {
	log *zerolog.Logger
}

func AsyncLogLevel(l zerolog.Level) asynq.LogLevel {
	switch l {
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

func NewAsynqLogger(log *zerolog.Logger) *asynqLogger {
	return &asynqLogger{log: log}
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

var _ tgbotapi.BotLogger = (*botLogger)(nil)

type botLogger struct {
	log *zerolog.Logger
}

func NewBotLogger(log *zerolog.Logger) *botLogger {
	return &botLogger{log: log}
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
