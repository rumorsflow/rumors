package front

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"golang.org/x/exp/slog"
)

var uiBuiltIn = true

type Front struct {
	Logger      *slog.Logger
	FeedRepo    repository.ReadRepository[*entity.Feed]
	ArticleRepo repository.ReadRepository[*entity.Article]
}

func (front *Front) Register(mux *wool.Wool) {
	mux.Group("/api/v1", func(w *wool.Wool) {
		feedActions := &FeedActions{FeedRepo: front.FeedRepo}
		articleActions := &ArticleActions{ArticleRepo: front.ArticleRepo, FeedRepo: front.FeedRepo}

		w.GET("/feeds", feedActions.List)
		w.GET("/articles", articleActions.List)
	})

	front.Logger.WithGroup("api").WithGroup("v1").Info("frontend V1 APIs registered")

	if uiBuiltIn {
		mux.UI("", assetFS())

		front.Logger.WithGroup("ui").Info("frontend UI registered")
	}
}
