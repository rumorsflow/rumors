package http

import (
	"fmt"
	"github.com/gowool/wool"
	"golang.org/x/exp/slog"
	"time"
)

func AfterServe(cfg *LogReqConfig, log *slog.Logger) wool.AfterServe {
	return func(c wool.Ctx, start, end time.Time, err error) {
		method := c.Req().Method
		endpoint := c.Req().URL.Path
		status := c.Res().Status()

		if !cfg.isOK(fmt.Sprintf("%d", status), method, endpoint) {
			return
		}

		host := c.Req().URL.Host
		if host == "" {
			host = c.Req().Host
		}

		args := []any{
			"status", status,
			"method", method,
			"host", host,
			"path", endpoint,
			"query", c.Req().URL.RawQuery,
			"ip", c.Req().RemoteAddr,
			"user-agent", c.Req().UserAgent(),
			"duration", end.Sub(start),
			"latency", fmt.Sprintf("%s", end.Sub(start)),
			"referer", c.Req().Referer(),
		}

		if err != nil {
			args = append(args, nil, nil)
			copy(args[2:], args)
			args[0] = "err"
			args[1] = err
			log.Error(c.Req().URL.String(), args...)
		} else if status >= 400 {
			log.Warn(c.Req().URL.String(), args...)
		} else {
			log.Info(c.Req().URL.String(), args...)
		}
	}
}
