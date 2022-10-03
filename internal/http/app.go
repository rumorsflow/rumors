package http

import (
	"context"
	"github.com/iagapie/rumors/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
	"net"
	"net/http"
	"time"
)

type App struct {
	cfg    config.ServerConfig
	log    zerolog.Logger
	e      *echo.Echo
	server *http.Server
}

func NewApp(debug bool, cfg config.ServerConfig, validator echo.Validator) *App {
	l := log.Logger.With().Str("context", "http").Logger()

	e := echo.New()
	e.Debug = debug
	e.HideBanner = true
	e.Validator = validator
	e.Logger = lecho.From(l)

	e.Use(middleware.Recover(), middleware.RemoveTrailingSlash())

	return &App{
		cfg: cfg,
		log: l,
		e:   e,
		server: &http.Server{
			Handler: e,
		},
	}
}

func (a *App) Echo() *echo.Echo {
	return a.e
}

func (a *App) Start() {
	a.log.Info().Msg("Start HTTP Server")

	listener, err := net.Listen(a.cfg.Network, a.cfg.Address)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to start HTTP Server")
	}

	if err = a.server.Serve(listener); err != http.ErrServerClosed {
		a.log.Error().Err(err).Msg("Failed to start HTTP Server")
	}
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.log.Info().Msg("Shutdown HTTP Server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Warn().Err(err).Msg("Failed to shut down server within given timeout.")
	}
}
