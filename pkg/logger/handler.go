package logger

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"io"
	"runtime"
	"strconv"
	"strings"
)

var (
	_ slog.Handler  = (*ConsoleHandler)(nil)
	_ HandlerSyncer = (*handlerSyncer)(nil)
)

const timeLayout = "2006-01-02T15:04:05.000Z0700"

type HandlerOptions struct {
	slog.HandlerOptions
}

func (opts HandlerOptions) NewConsoleHandler(w io.Writer) *ConsoleHandler {
	return NewConsoleHandler(opts.HandlerOptions, w)
}

func (opts HandlerOptions) NewHandler(w io.Writer, encoding string) slog.Handler {
	switch strings.ToLower(encoding) {
	case "console":
		return opts.NewConsoleHandler(w)
	case "text":
		return opts.NewTextHandler(w)
	default:
		return opts.NewJSONHandler(w)
	}
}

type ConsoleHandler struct {
	opts   slog.HandlerOptions
	global []slog.Attr
	groups []string
	w      io.Writer
}

func NewConsoleHandler(opts slog.HandlerOptions, w io.Writer, attrs ...slog.Attr) *ConsoleHandler {
	return &ConsoleHandler{opts: opts, w: w, global: attrs}
}

func (h *ConsoleHandler) clone() *ConsoleHandler {
	return &ConsoleHandler{
		global: slices.Clip(h.global),
		groups: slices.Clip(h.groups),
		opts:   h.opts,
		w:      h.w,
	}
}

func (h *ConsoleHandler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := h.clone()
	c.global = append(c.global, attrs...)
	return c
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	c := h.clone()
	c.groups = append(c.groups, name)
	return c
}

func (h *ConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	if !r.Time.IsZero() {
		_, _ = buf.WriteString(spaces(r.Time.Format(timeLayout), 31))
	}

	_, _ = buf.WriteString(coloredLevel(r.Level))

	if len(h.groups) > 0 {
		_, _ = buf.WriteString(coloredGroup(strings.Join(h.groups, ".")))
	}

	if r.Message != "" {
		_, _ = buf.WriteString(spaces(r.Message, 24))
	}

	attrs, sep := h.attrs(r)
	attrs += h.addSource(r, sep)

	if attrs != "" {
		_, _ = buf.WriteString(" {")
		_, _ = buf.WriteString(attrs)
		_ = buf.WriteByte('}')
	}

	if err := buf.WriteByte('\n'); err != nil {
		return err
	}

	_, err := h.w.Write(buf.Bytes())
	return err
}

func (h *ConsoleHandler) attrs(r slog.Record) (string, string) {
	total := len(h.global) + r.NumAttrs()

	if total == 0 {
		return "", ""
	}

	sep := ""

	var buf bytes.Buffer
	var fn func(a slog.Attr) bool

	fn = func(a slog.Attr) bool {
		total--

		v := a.Value.Resolve()

		_, _ = buf.WriteString(sep)
		_ = buf.WriteByte('"')
		_, _ = buf.WriteString(a.Key)
		_, _ = buf.WriteString("\": ")

		if v.Kind() == slog.KindGroup {
			sep = ""
			_ = buf.WriteByte('{')
			for _, aa := range v.Group() {
				_, _ = buf.WriteString(sep)
				fn(aa)
			}
			_ = buf.WriteByte('}')
		} else {
			sep = ", "
			if err, ok := v.Any().(error); ok {
				_, _ = buf.WriteString(err.Error())
			} else {
				b, _ := json.MarshalWithOption(v.Any())
				_, _ = buf.Write(b)
			}
		}
		return true
	}

	for _, attr := range h.global {
		fn(attr)
	}

	r.Attrs(fn)

	if h.opts.AddSource {
		f := frame(r)
		if f.File != "" {
			_, _ = buf.WriteString(sep)
			_ = buf.WriteByte('"')
			_, _ = buf.WriteString(slog.SourceKey)
			_, _ = buf.WriteString("\": ")
			_ = buf.WriteByte('"')
			_, _ = buf.WriteString(f.File)
			_ = buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(f.Line))
			_ = buf.WriteByte('"')
		}
	}

	return buf.String(), sep
}

func (h *ConsoleHandler) addSource(r slog.Record, sep string) string {
	if !h.opts.AddSource {
		return ""
	}

	f := frame(r)

	if f.File != "" {
		var buf bytes.Buffer

		_, _ = buf.WriteString(sep)
		_ = buf.WriteByte('"')
		_, _ = buf.WriteString(slog.SourceKey)
		_, _ = buf.WriteString("\": ")
		_ = buf.WriteByte('"')
		_, _ = buf.WriteString(f.File)
		_ = buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(f.Line))
		_ = buf.WriteByte('"')

		return buf.String()
	}

	return ""
}

func frame(r slog.Record) runtime.Frame {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return f
}

func spaces(str string, min int) string {
	if len(str) < min {
		return str + strings.Repeat(" ", min-len(str)) + " "
	}
	return str + " "
}

func coloredLevel(level slog.Level) string {
	str := spaces(level.String(), 7)

	switch level {
	case slog.LevelInfo:
		return color.HiCyanString(str)
	case slog.LevelWarn:
		return color.HiYellowString(str)
	case slog.LevelError:
		return color.HiRedString(str)
	default:
		return color.HiWhiteString(str)
	}
}

func coloredGroup(group string) string {
	return color.HiGreenString(spaces(group, 16))
}

type HandlerSyncer interface {
	Sync() error
}

type handlerSyncer struct {
	slog.Handler
	syncer WriteSyncer
}

func (h *handlerSyncer) Sync() error {
	return h.syncer.Sync()
}
