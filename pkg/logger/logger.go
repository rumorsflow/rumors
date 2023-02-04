package logger

import (
	"context"
	"golang.org/x/exp/slog"
	"io/fs"
	"strings"
)

var syncer WriteSyncer

func Init(cfg *Config) {
	if cfg == nil {
		panic("logger config is nil")
	}

	cfg.Init()

	var err error
	if syncer, err = cfg.openSinks(); err != nil {
		panic(err)
	}

	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	encoding := strings.ToLower(cfg.Encoding)

	if cfg.Development {
		level = slog.LevelDebug
		encoding = "console"
	}

	var handler slog.Handler

	switch encoding {
	case "json":
		opts := slog.HandlerOptions{Level: level, AddSource: cfg.AddSource}
		handler = opts.NewJSONHandler(syncer)
	case "console":
		handler = NewConsoleHandler(syncer, level, cfg.AddSource)
	default:
		opts := slog.HandlerOptions{Level: level, AddSource: cfg.AddSource}
		handler = opts.NewTextHandler(syncer)
	}

	if len(cfg.Attrs) > 0 {
		attrs := make([]slog.Attr, 0, len(cfg.Attrs))
		for key, value := range cfg.Attrs {
			attrs = append(attrs, slog.Any(key, value))
		}

		handler = handler.WithAttrs(attrs)
	}

	slog.SetDefault(slog.New(handler))
}

func Sync() error {
	if err := syncer.Sync(); err != nil {
		if pe, ok := err.(*fs.PathError); ok && (strings.Contains(pe.Path, "stderr") || strings.Contains(pe.Path, "stdout")) {
			return nil
		}
		return err
	}
	return nil
}

func IsDebug() bool {
	return Enabled(slog.LevelDebug)
}

func Enabled(level slog.Level) bool {
	return slog.Default().Enabled(level)
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

func Error(msg string, err error, args ...any) {
	slog.Error(msg, err, args...)
}

func Log(level slog.Level, msg string, args ...any) {
	slog.Log(level, msg, args...)
}

func LogAttrs(level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(level, msg, attrs...)
}

func With(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}

func WithGroup(name string) *slog.Logger {
	return slog.Default().WithGroup(name)
}

func WithContext(ctx context.Context) *slog.Logger {
	return slog.Default().WithContext(ctx)
}
