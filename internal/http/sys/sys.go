package sys

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"golang.org/x/exp/slog"
	"strings"
)

var uiBuiltIn = true

type Sys struct {
	Logger      *slog.Logger
	CfgJWT      *jwt.Config
	AuthService AuthService
	FeedRepo    repository.ReadWriteRepository[*entity.Feed]
	ChatRepo    repository.ReadWriteRepository[*entity.Chat]
	JobRepo     repository.ReadWriteRepository[*entity.Job]
	ArticleRepo repository.ReadWriteRepository[*entity.Article]
}

func (s *Sys) Register(mux *wool.Wool) {
	mux.Group("/sys", func(sys *wool.Wool) {
		sys.Group("/api", func(w *wool.Wool) {
			w.Group("/auth", func(a *wool.Wool) {
				auth := NewAuthActions(s.AuthService)

				a.POST("/sign-in", auth.SignIn)
				a.POST("/refresh", auth.Refresh)

				a.Use(JWTMiddleware(s.CfgJWT, false))
				a.POST("/otp", auth.OTP)
			})

			w.Use(JWTMiddleware(s.CfgJWT, true))
			w.CRUD("/feeds", NewFeedCRUD(s.FeedRepo, s.FeedRepo))
			w.CRUD("/chats", NewChatCRUD(s.ChatRepo, s.ChatRepo))
			w.CRUD("/jobs", NewJobCRUD(s.JobRepo, s.JobRepo))
			w.CRUD("/articles", NewArticleActions(s.ArticleRepo, s.ArticleRepo))
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
