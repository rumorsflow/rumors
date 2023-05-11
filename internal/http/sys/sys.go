package sys

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
)

var uiBuiltIn = true

type Sys struct {
	Logger         *slog.Logger
	CfgJWT         *jwt.Config
	SSE            *SSE
	AuthActions    *AuthActions
	QueueActions   *QueueActions
	ArticleActions *ArticleActions
	SiteCRUD       action.CRUD
	ChatCRUD       action.CRUD
	JobCRUD        action.CRUD
	DirUI          string
}

func (s *Sys) Register(mux *wool.Wool) {
	mux.Group("/sys", func(sys *wool.Wool) {
		sys.Group("/api", func(w *wool.Wool) {
			w.Group("/auth", func(a *wool.Wool) {
				a.POST("/sign-in", s.AuthActions.SignIn)
				a.POST("/refresh", s.AuthActions.Refresh)

				a.Group("", func(g *wool.Wool) {
					g.Use(JWTMiddleware(s.CfgJWT, false))
					g.POST("/otp", s.AuthActions.OTP)
				})

				a.Group("", func(g *wool.Wool) {
					g.Use(JWTMiddleware(s.CfgJWT, true))
					g.POST("/sse", s.SSE.Auth)
				})
			})

			w.Group("", func(sw *wool.Wool) {
				sw.Use(s.SSE.Middleware)
				sw.GET("/realtime", s.SSE.Handler)
			})

			w.Use(JWTMiddleware(s.CfgJWT, true))

			w.CRUD("/articles", s.ArticleActions)
			w.CRUD("/sites", s.SiteCRUD)
			w.CRUD("/chats", s.ChatCRUD)
			w.CRUD("/jobs", s.JobCRUD)

			w.Group("/queues", func(q *wool.Wool) {
				q.DELETE("/:"+QNameParam, s.QueueActions.Delete)
				q.POST("/:"+QNameParam+"/pause", s.QueueActions.Pause)
				q.POST("/:"+QNameParam+"/resume", s.QueueActions.Resume)
			})
		})

		s.Logger.WithGroup("api").Info("system APIs registered")

		if uiBuiltIn || s.DirUI != "" {
			sys.Use(func(next wool.Handler) wool.Handler {
				return func(c wool.Ctx) error {
					path := c.Req().URL.Path
					c.Req().URL.Path = strings.TrimPrefix(c.Req().URL.Path, "/sys")
					err := next(c)
					c.Req().URL.Path = path
					return err
				}
			})

			if s.DirUI != "" {
				sys.UI("", http.Dir(s.DirUI))
				s.Logger.WithGroup("ui").Info("system UI registered")
			} else {
				sys.UI("", sysAssetFS())
				s.Logger.WithGroup("ui").Info("system UI registered")
			}
		}
	})
}

func (s *Sys) Listen(done <-chan struct{}) {
	go s.SSE.Listen(done)
}
