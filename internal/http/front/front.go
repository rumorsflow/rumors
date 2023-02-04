package front

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"golang.org/x/exp/slog"
)

var uiBuiltIn = true

type Front struct {
	Logger           *slog.Logger
	FeedRepo         repository.ReadRepository[*entity.Feed]
	ArticleRepo      repository.ReadRepository[*entity.Article]
	FeedNoFilters    []string
	ArticleNoFilters []string
}

func (front *Front) Register(mux *wool.Wool) {
	mux.Group("/api/v1", func(w *wool.Wool) {
		w.CRUD("/feeds", NewFeedActions(front.FeedRepo, front.FeedNoFilters))
		w.CRUD("/articles", NewArticleActions(front.ArticleRepo, front.ArticleNoFilters))
	})

	front.Logger.WithGroup("api").WithGroup("v1").Info("frontend V1 APIs registered")

	if uiBuiltIn {
		mux.UI("", assetFS())

		front.Logger.WithGroup("ui").Info("frontend UI registered")
	}
}
