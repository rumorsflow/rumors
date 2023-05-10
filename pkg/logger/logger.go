package logger

import (
	"github.com/rumorsflow/rumors/v2/pkg/util"
	"golang.org/x/exp/slog"
	"strings"
)

var _ Logger = (*Log)(nil)

type Logger interface {
	NamedLogger(name string) *slog.Logger
}

type Log struct {
	base     *slog.Logger
	channels ChannelConfig
}

func NewLogger(channels ChannelConfig, base *slog.Logger) *Log {
	return &Log{
		channels: channels,
		base:     base,
	}
}

func (l *Log) NamedLogger(name string) *slog.Logger {
	if cfg, ok := l.channels.Channels[name]; ok {
		return util.Must(cfg.Logger()).WithGroup(name)
	}

	return l.base.WithGroup(name)
}

func ToLeveler(level string) slog.Leveler {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func ToAttrs(data map[string]any) []slog.Attr {
	var attrs []slog.Attr
	for key, value := range data {
		attrs = append(attrs, slog.Any(key, value))
	}
	return attrs
}

func (cfg *Config) Logger() (*slog.Logger, error) {
	syncer, err := cfg.OpenSinks()
	if err != nil {
		return nil, err
	}

	handler := cfg.Opts().NewHandler(syncer, cfg.Encoding)

	attrs := ToAttrs(cfg.Attrs)
	if len(attrs) > 0 {
		handler = handler.WithAttrs(attrs)
	}

	return slog.New(&handlerSyncer{
		Handler: handler,
		syncer:  syncer,
	}), nil
}
