package gzip

import (
	"github.com/NYTimes/gziphandler"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"net/http"
)

const (
	RootPluginName = "http"
	PluginName     = "gzip"
)

type Plugin struct{}

func (*Plugin) Init(cfg config.Configurer) error {
	const op = errors.Op("gzip plugin init")

	if !cfg.Has(RootPluginName) {
		return errors.E(op, errors.Disabled)
	}

	return nil
}

// Name returns user-friendly plugin name
func (*Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Handle(next http.Handler) http.Handler {
	return gziphandler.GzipHandler(next)
}
