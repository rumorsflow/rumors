package front

import (
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type ArticleActions struct {
	*action.ListAction[*entity.Article]
}

func NewArticleActions(read repository.ReadRepository[*entity.Article], noFilters []string) *ArticleActions {
	return &ArticleActions{
		ListAction: &action.ListAction[*entity.Article]{
			ReadRepository: read,
			NoFilters:      noFilters,
		},
	}
}
