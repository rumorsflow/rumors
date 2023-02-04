package front

import (
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type FeedActions struct {
	*action.ListAction[*entity.Feed]
}

func NewFeedActions(read repository.ReadRepository[*entity.Feed], noFilters []string) *FeedActions {
	return &FeedActions{
		ListAction: &action.ListAction[*entity.Feed]{
			ReadRepository: read,
			NoFilters:      noFilters,
		},
	}
}
