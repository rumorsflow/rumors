package sys

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"golang.org/x/exp/slog"
	"strings"
)

var uiBuiltIn = true

type Sys struct {
	Logger         *slog.Logger
	CfgJWT         *jwt.Config
	AuthActions    *AuthActions
	ArticleActions *ArticleActions
	FeedCRUD       action.CRUD
	ChatCRUD       action.CRUD
	JobCRUD        action.CRUD
}

func (s *Sys) Register(mux *wool.Wool) {
	mux.Group("/sys", func(sys *wool.Wool) {
		sys.Group("/api", func(w *wool.Wool) {
			w.Group("/auth", func(a *wool.Wool) {
				a.POST("/sign-in", s.AuthActions.SignIn)
				a.POST("/refresh", s.AuthActions.Refresh)

				a.Use(JWTMiddleware(s.CfgJWT, false))
				a.POST("/otp", s.AuthActions.OTP)
			})

			w.Use(JWTMiddleware(s.CfgJWT, true))
			w.CRUD("/articles", s.ArticleActions)
			w.CRUD("/feeds", s.FeedCRUD)
			w.CRUD("/chats", s.ChatCRUD)
			w.CRUD("/jobs", s.JobCRUD)
		})

		s.Logger.WithGroup("api").Info("system APIs registered")

		if uiBuiltIn {
			sys.Use(func(next wool.Handler) wool.Handler {
				return func(c wool.Ctx) error {
					path := c.Req().URL.Path
					c.Req().URL.Path = strings.TrimPrefix(c.Req().URL.Path, "/sys")
					err := next(c)
					c.Req().URL.Path = path
					return err
				}
			})
			sys.UI("", sysAssetFS())

			s.Logger.WithGroup("ui").Info("system UI registered")
		}
	})
}
