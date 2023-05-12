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
	attrs    map[string]any
	base     *slog.Logger
	channels ChannelConfig
}

func NewLogger(attrs map[string]any, channels ChannelConfig, base *slog.Logger) *Log {
	return &Log{
		attrs:    attrs,
		channels: channels,
		base:     base,
	}
}

func (l *Log) NamedLogger(name string) *slog.Logger {
	if cfg, ok := l.channels.Channels[name]; ok {
		return util.Must(cfg.Logger(l.attrs)).WithGroup(name)
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
	attrs := make([]slog.Attr, 0, len(data))
	for key, value := range data {
		attrs = append(attrs, slog.Any(key, value))
	}
	return attrs
}

func (cfg *Config) Logger(attrs map[string]any) (*slog.Logger, error) {
	syncer, err := cfg.OpenSinks()
	if err != nil {
		return nil, err
	}

	handler := cfg.Opts().NewHandler(syncer, cfg.Encoding)
	handler = handler.WithAttrs(append(ToAttrs(cfg.Attrs), ToAttrs(attrs)...))

	return slog.New(&handlerSyncer{
		Handler: handler,
		syncer:  syncer,
	}), nil
}
