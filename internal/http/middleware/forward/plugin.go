package forward

import (
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"net/http"
	"path/filepath"
)

const (
	RootPluginName = "http"
	PluginName     = "forward"
)

type Plugin struct{}

func (*Plugin) Init(cfg config.Configurer) error {
	const op = errors.Op("forward plugin init")

	if !cfg.Has(RootPluginName) {
		return errors.E(op, errors.Disabled)
	}

	return nil
}

// Name returns user-friendly plugin name
func (*Plugin) Name() string {
	return PluginName
}

func (*Plugin) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &wrapperResponseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		if ww.statusCode != 404 {
			return
		}

		if ext := filepath.Ext(r.URL.Path); ext == "" || ext == ".html" {
			r.URL.Path = "/"
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(404)
			_, _ = w.Write(ww.data)
		}
	})
}

type wrapperResponseWriter struct {
	http.ResponseWriter
	statusCode int
	data       []byte
}

func (w *wrapperResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	if statusCode != 404 {
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *wrapperResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 404 {
		w.data = data
		return 0, nil
	}
	return w.ResponseWriter.Write(data)
}
