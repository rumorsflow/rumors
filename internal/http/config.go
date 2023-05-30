package http

import (
	"github.com/gowool/middleware/cfipcountry"
	"github.com/gowool/middleware/cors"
	"github.com/gowool/middleware/gzip"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/swagger"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"regexp"
	"time"
)

type UIConfig struct {
	FrontPath string `mapstructure:"front"`
	SysPath   string `mapstructure:"sys"`
}

type LogReqConfig struct {
	ExcludeRegexStatus   string `mapstructure:"exclude_status"`
	ExcludeRegexMethod   string `mapstructure:"exclude_method"`
	ExcludeRegexEndpoint string `mapstructure:"exclude_endpoint"`

	rxStatus   *regexp.Regexp
	rxMethod   *regexp.Regexp
	rxEndpoint *regexp.Regexp
}

func (cfg *LogReqConfig) init() {
	if cfg.ExcludeRegexEndpoint == "" {
		cfg.ExcludeRegexEndpoint = "^/(metrics|favicon.ico)"
	}

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

func (cfg *LogReqConfig) isOK(status, method, endpoint string) bool {
	return (cfg.rxStatus == nil || !cfg.rxStatus.MatchString(status)) &&
		(cfg.rxMethod == nil || !cfg.rxMethod.MatchString(method)) &&
		(cfg.rxEndpoint == nil || !cfg.rxEndpoint.MatchString(endpoint))
}

type MiddlewareConfig struct {
	Compress    *gzip.Config        `mapstructure:"compress"`
	CORS        *cors.Config        `mapstructure:"cors"`
	Metrics     *prometheus.Config  `mapstructure:"metrics"`
	CfIPCountry *cfipcountry.Config `mapstructure:"cfipcountry"`
}

type SwaggerConfig struct {
	Enabled bool            `mapstructure:"enabled"`
	Front   *swagger.Config `mapstructure:"front"`
	Sys     *swagger.Config `mapstructure:"sys"`
}

type Config struct {
	UI         *UIConfig         `mapstructure:"ui"`
	LogReq     *LogReqConfig     `mapstructure:"log_request"`
	Middleware *MiddlewareConfig `mapstructure:"middleware"`
	Swagger    *SwaggerConfig    `mapstructure:"swagger"`
	JWT        *jwt.Config       `mapstructure:"jwt"`
	SSE        *sse.Config       `mapstructure:"-"`
}

func (cfg *Config) init(version string) {
	if cfg.UI == nil {
		cfg.UI = &UIConfig{}
	}

	if cfg.LogReq == nil {
		cfg.LogReq = &LogReqConfig{}
	}
	cfg.LogReq.init()

	cfg.SSE = &sse.Config{
		ClientIdle: time.Hour,
		Metrics: &sse.MetricsConfig{
			Enabled: false,
			Version: version,
		},
	}
	cfg.SSE.Init()

	if cfg.Middleware == nil {
		cfg.Middleware = &MiddlewareConfig{}
	}

	if cfg.Middleware.Compress == nil {
		cfg.Middleware.Compress = &gzip.Config{}
	}
	cfg.Middleware.Compress.Init()

	if cfg.Middleware.CORS == nil {
		cfg.Middleware.CORS = &cors.Config{}
	}

	if cfg.Middleware.Metrics != nil {
		cfg.SSE.Metrics.Enabled = true
		cfg.Middleware.Metrics.Version = version
		cfg.Middleware.Metrics.Init()
	}

	if cfg.Swagger == nil {
		cfg.Swagger = &SwaggerConfig{}
	}

	if cfg.Swagger.Enabled {
		if cfg.Swagger.Front == nil {
			cfg.Swagger.Front = &swagger.Config{}
		}
		cfg.Swagger.Front.InstanceName = "front"
		cfg.Swagger.Front.Init()

		if cfg.Swagger.Sys == nil {
			cfg.Swagger.Sys = &swagger.Config{}
		}
		cfg.Swagger.Sys.InstanceName = "sys"
		cfg.Swagger.Sys.Init()
	}

	if cfg.JWT == nil {
		cfg.JWT = &jwt.Config{}
	}
	cfg.JWT.Init()
}
