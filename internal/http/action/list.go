package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"net/http"
)

var _ wool.List = (*ListAction[repository.Entity])(nil)

type ListResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
	Index int64 `json:"index"`
	Size  int64 `json:"size"`
}

type ListAction[Entity repository.Entity] struct {
	ReadRepository repository.ReadRepository[Entity]
	NoFilters      []string
}

func (a *ListAction[Entity]) List(ctx wool.Ctx) error {
	criteria := db.BuildCriteria(ctx.Req().URL.RawQuery, a.NoFilters...)
	if criteria.Index == nil {
		criteria.SetIndex(0)
	}
	if criteria.Size == nil {
		criteria.SetSize(20)
	}

	total, err := a.ReadRepository.Count(ctx.Req().Context(), criteria.Filter)
	if err != nil {
		return err
	}

	response := ListResponse{
		Total: total,
		Index: *criteria.Index,
		Size:  *criteria.Size,
	}

	if total > 0 {
		response.Data, err = a.ReadRepository.Find(ctx.Req().Context(), criteria)
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
