package http

import (
	"context"
	"github.com/iagapie/rumors/internal/config"
	"github.com/iagapie/rumors/pkg/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
	"net"
	"net/http"
	"time"
)

type App struct {
	cfg    config.ServerConfig
	log    *zerolog.Logger
	e      *echo.Echo
	server *http.Server
}

func NewApp(debug bool, cfg config.ServerConfig, log *zerolog.Logger) *App {
	e := echo.New()
	e.Debug = debug
	e.HideBanner = true
	e.Validator = validate.New()
	e.Logger = lecho.From(*log)

	e.Use(middleware.Recover(), middleware.RemoveTrailingSlash())

	return &App{
		cfg: cfg,
		log: log,
		e:   e,
		server: &http.Server{
			Handler: e,
		},
	}
}

func (a *App) Echo() *echo.Echo {
	return a.e
}

func (a *App) Start() error {
	a.log.Info().Msg("Start HTTP Server")

	listener, err := net.Listen(a.cfg.Network, a.cfg.Address)
	if err != nil {
		return err
	}

	if err = a.server.Serve(listener); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.log.Info().Msg("Shutdown HTTP Server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Warn().Err(err).Msg("Failed to shut down server within given timeout.")
	}
}
