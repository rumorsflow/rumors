package auth

import (
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/rumors/internal/api/util"
	"go.uber.org/zap"
	"net/http"
)

const PluginName = "/api/v1/auth"

type Plugin struct {
	log *zap.Logger
}

func (p *Plugin) Init(log *zap.Logger) error {
	p.log = log
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Register(mux *flow.Mux) {
	mux.HandleFunc(PluginName+"/session", p.session, http.MethodPost)
	mux.HandleFunc(PluginName+"/otp", p.otp, http.MethodPost)
	mux.HandleFunc(PluginName+"/refresh", p.refresh, http.MethodPost)
}

func (p *Plugin) session(w http.ResponseWriter, r *http.Request) {
	var dto SessionRequest
	util.Bind(r, &dto)
}

func (p *Plugin) otp(w http.ResponseWriter, r *http.Request) {
	var dto OtpRequest
	util.Bind(r, &dto)
}

func (p *Plugin) refresh(w http.ResponseWriter, r *http.Request) {
	var dto RefreshTokenRequest
	util.Bind(r, &dto)
}
