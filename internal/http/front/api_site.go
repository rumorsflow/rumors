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

type SiteActions struct {
	SiteRepo repository.ReadRepository[*entity.Site]
}

func (a *SiteActions) List(c wool.Ctx) error {
	filter := bson.M{"enabled": true}

	total, err := a.SiteRepo.Count(c.Req().Context(), filter)
	if err != nil {
		return err
	}

	query := c.Req().URL.Query()

	criteria := &repository.Criteria{Sort: bson.D{{Key: "domain", Value: 1}}, Filter: filter}
	criteria.SetIndex(cast.ToInt64(query.Get(db.QueryIndex)))
	criteria.SetSize(cast.ToInt64(query.Get(db.QuerySize)))

	response := action.ListResponse{
		Total: total,
		Index: *criteria.Index,
		Size:  *criteria.Size,
	}

	if total > 0 {
		response.Data, err = a.SiteRepo.Find(c.Req().Context(), criteria)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, response)
}
