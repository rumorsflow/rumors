package jwt

import (
	"context"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/pkg/jwt"
	"github.com/rumorsflow/rumors/internal/services/token"
	"net/http"
	"strings"
)

const (
	RootPluginName = "http"
	PluginName     = "jwt"

	header = "Authorization"
	prefix = "Bearer "
)

type Plugin struct {
	service token.Service
}

type (
	ctxClaimsKey struct{}
	ctxUserKey   struct{}
)

func CtxClaims(ctx context.Context) (claims jwt.UserClaims, ok bool) {
	claims, ok = ctx.Value(ctxClaimsKey{}).(jwt.UserClaims)
	return
}

func CtxUser(ctx context.Context) (user models.User, ok bool) {
	user, ok = ctx.Value(ctxUserKey{}).(models.User)
	return
}

func CtxUserId(ctx context.Context) string {
	user, _ := CtxUser(ctx)
	return user.Id
}

func (p *Plugin) Init(cfg config.Configurer, service token.Service) error {
	const op = errors.Op("http jwt plugin init")

	if !cfg.Has(RootPluginName) {
		return errors.E(op, errors.Disabled)
	}

	p.service = service

	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t := p.getToken(r); t != "" {
			claims, user, err := p.service.ValidateJWT(r.Context(), t)
			if err != nil {
				util.Unauthorized(err)
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxClaimsKey{}, claims)
			ctx = context.WithValue(ctx, ctxUserKey{}, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) getToken(r *http.Request) string {
	for _, value := range r.Header.Values(header) {
		if len(value) > len(prefix) && strings.EqualFold(value[:len(prefix)], prefix) {
			return value[len(prefix):]
		}
	}
	return ""
}
