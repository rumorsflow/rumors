package http

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"github.com/rumorsflow/rumors/internal/pkg/validate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

const PluginName = "http"

type Plugin struct {
	cfg *Config
	e   *echo.Echo
	log *zap.Logger
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	const op = errors.Op("http plugin init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, errors.Init, err)
	}
	p.cfg.InitDefault()

	p.log = log
	p.e = echo.New()
	p.e.Debug = zapcore.LevelOf(log.Core()) == zapcore.DebugLevel
	p.e.Logger.SetOutput(io.Discard)
	p.e.StdLogger.SetOutput(io.Discard)
	p.e.ListenerNetwork = p.cfg.Network
	p.e.HideBanner = true
	p.e.Validator = validate.New()
	p.e.Use(echozap.ZapLogger(log), middleware.Recover(), middleware.RemoveTrailingSlash())
	return nil
}

func (p *Plugin) Serve() chan error {
	return make(chan error, 1)
}

func (p *Plugin) Stop() error {
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
