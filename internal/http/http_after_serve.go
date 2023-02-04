package http

import (
	"fmt"
	"github.com/gowool/wool"
	"regexp"
	"time"
)

type AfterServeConfig struct {
	ExcludeRegexStatus   string `mapstructure:"exclude_status"`
	ExcludeRegexMethod   string `mapstructure:"exclude_method"`
	ExcludeRegexEndpoint string `mapstructure:"exclude_endpoint"`

	rxStatus   *regexp.Regexp
	rxMethod   *regexp.Regexp
	rxEndpoint *regexp.Regexp
}

func (cfg *AfterServeConfig) Init() {
	if cfg.ExcludeRegexStatus != "" {
		cfg.rxStatus, _ = regexp.Compile(cfg.ExcludeRegexStatus)
	}
	if cfg.ExcludeRegexMethod != "" {
		cfg.rxMethod, _ = regexp.Compile(cfg.ExcludeRegexMethod)
	}
	if cfg.ExcludeRegexEndpoint != "" {
		cfg.rxEndpoint, _ = regexp.Compile(cfg.ExcludeRegexEndpoint)
	}
}

func (cfg *AfterServeConfig) isOK(status, method, endpoint string) bool {
	return (cfg.rxStatus == nil || !cfg.rxStatus.MatchString(status)) &&
		(cfg.rxMethod == nil || !cfg.rxMethod.MatchString(method)) &&
		(cfg.rxEndpoint == nil || !cfg.rxEndpoint.MatchString(endpoint))
}

func AfterServe(cfg *AfterServeConfig) wool.AfterServe {
	if cfg == nil {
		cfg = &AfterServeConfig{}
	}
	cfg.Init()

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
		}

		if err != nil {
			wool.Logger().Error(c.Req().URL.String(), err, args...)
		} else if status >= 400 {
			wool.Logger().Warn(c.Req().URL.String(), args...)
		} else {
			wool.Logger().Info(c.Req().URL.String(), args...)
		}
	}
}
