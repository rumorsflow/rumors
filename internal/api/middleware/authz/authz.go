package authz

import (
	"errors"
	"github.com/rumorsflow/rumors/internal/api/middleware/jwt"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/rumors/internal/models"
	"net/http"
)

var (
	ErrJWTMissing = errors.New("missing or malformed jwt")
	ErrNoOTP      = errors.New("2FA is required")
	ErrNoRole     = errors.New("access denied")
)

func Is(next http.HandlerFunc, roles ...models.Role) http.HandlerFunc {
	return middleware(false, roles...)(next)
}

func Has(next http.HandlerFunc, roles ...models.Role) http.HandlerFunc {
	return middleware(true, roles...)(next)
}

func IsAdmin(next http.HandlerFunc, roles ...models.Role) http.HandlerFunc {
	return Has(next, append([]models.Role{models.AdminRole}, roles...)...)
}

func middleware(full bool, roles ...models.Role) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok1 := jwt.CtxClaims(r.Context())
			user, ok2 := jwt.CtxUser(r.Context())

			if !ok1 || !ok2 {
				util.BadRequest(ErrJWTMissing)
			}

			if otp, _ := claims.Metadata.(bool); full && !otp {
				util.Forbidden(ErrNoOTP)
			}

			if len(roles) > 0 && !user.IsGranted(roles...) {
				util.Forbidden(ErrNoRole)
			}

			next(w, r)
		}
	}
}
