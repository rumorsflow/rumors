package front

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type FeedActions struct {
	*action.ListAction[*entity.Feed, any]
}

func NewFeedActions(read repository.ReadRepository[*entity.Feed]) *FeedActions {
	return &FeedActions{
		ListAction: &action.ListAction[*entity.Feed, any]{
			ReadRepository: read,
			CriteriaBuilder: func(c wool.Ctx) *repository.Criteria {
				criteria := action.DefaultCriteriaBuilder(c, "enabled")
				if ff, ok := criteria.Filter.(bson.M); ok {
					ff["enabled"] = true
				}
				return criteria
			},
		},
	}
}
