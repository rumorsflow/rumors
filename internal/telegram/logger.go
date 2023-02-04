package telegram

import (
	"fmt"
	"golang.org/x/exp/slog"
)

type telegramLogger struct {
	logger *slog.Logger
}

func (tl *telegramLogger) Println(v ...any) {
	for _, item := range v {
		switch val := item.(type) {
		case error:
			tl.logger.Error("request error", val)
		case string:
			tl.logger.Info(val)
		default:
			tl.logger.Warn("unknown error", item)
		}
	}
}

func (tl *telegramLogger) Printf(format string, v ...any) {
	tl.logger.Info(fmt.Sprintf(format, v...))
}
