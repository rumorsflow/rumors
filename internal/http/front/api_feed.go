package front

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type FeedActions struct {
	FeedRepo repository.ReadRepository[*entity.Feed]
}

func (a *FeedActions) List(c wool.Ctx) error {
	filter := bson.M{"enabled": true}

	total, err := a.FeedRepo.Count(c.Req().Context(), filter)
	if err != nil {
		return err
	}

	query := c.Req().URL.Query()

	criteria := &repository.Criteria{
		Sort:   bson.D{{Key: "host", Value: 1}},
		Filter: filter,
	}
	criteria.SetIndex(cast.ToInt64(query.Get(db.QueryIndex)))
	criteria.SetSize(cast.ToInt64(query.Get(db.QuerySize)))

	response := action.ListResponse{
		Total: total,
		Index: *criteria.Index,
		Size:  *criteria.Size,
	}

	if total > 0 {
		response.Data, err = a.FeedRepo.Find(c.Req().Context(), criteria)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, response)
}
