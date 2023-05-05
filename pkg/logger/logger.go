package logger

import (
	"context"
	"golang.org/x/exp/slog"
	"io/fs"
	"strings"
)

func Init(cfg *Config) {
	if cfg == nil {
		panic("logger config is nil")
	}

	syncer, err := cfg.OpenSinks()
	if err != nil {
		panic(err)
	}

	handler := cfg.Opts().NewHandler(syncer, cfg.Encoding)

	attrs := ToAttrs(cfg.Attrs)
	if len(attrs) > 0 {
		handler = handler.WithAttrs(attrs)
	}

	slog.SetDefault(slog.New(&handlerSyncer{
		Handler: handler,
		syncer:  syncer,
	}))
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

func Sync() error {
	if syncer, ok := slog.Default().Handler().(HandlerSyncer); ok {
		if err := syncer.Sync(); err != nil {
			if pe, ok := err.(*fs.PathError); ok && (strings.Contains(pe.Path, "stderr") || strings.Contains(pe.Path, "stdout")) {
				return nil
			}
			return err
		}
	}
	return nil
}

func IsDebug() bool {
	return Enabled(slog.LevelDebug)
}

func Enabled(level slog.Level) bool {
	return slog.Default().Enabled(context.Background(), level)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
}

func With(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}

func WithGroup(name string) *slog.Logger {
	return slog.Default().WithGroup(name)
}
