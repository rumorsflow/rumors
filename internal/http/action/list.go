package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"net/http"
)

var _ wool.List = (*ListAction[repository.Entity, any])(nil)

var DefaultCriteriaBuilder = func(c wool.Ctx, excludeFields ...string) *repository.Criteria {
	criteria := db.BuildCriteria(c.Req().URL.RawQuery, excludeFields...)
	if criteria.Index == nil {
		criteria.SetIndex(0)
	}
	if criteria.Size == nil {
		criteria.SetSize(20)
	}

	return criteria
}

type ListResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
	Index int64 `json:"index"`
	Size  int64 `json:"size"`
}

type ListAction[Entity repository.Entity, DTO any] struct {
	ReadRepository  repository.ReadRepository[Entity]
	ResponseMapper  ResponseMapper[Entity, DTO]
	CriteriaBuilder func(c wool.Ctx) *repository.Criteria
}

func (a *ListAction[Entity, DTO]) List(c wool.Ctx) error {
	var criteria *repository.Criteria
	if a.CriteriaBuilder == nil {
		criteria = DefaultCriteriaBuilder(c)
	} else {
		criteria = a.CriteriaBuilder(c)
	}

	total, err := a.ReadRepository.Count(c.Req().Context(), criteria.Filter)
	if err != nil {
		return err
	}

	response := ListResponse{
		Total: total,
		Index: *criteria.Index,
		Size:  *criteria.Size,
	}

	if total > 0 {
		data, err := a.ReadRepository.Find(c.Req().Context(), criteria)
		if err != nil {
			return err
		}

		if a.ResponseMapper == nil {
			response.Data = data
		} else if len(data) > 0 {
			mapped := make([]DTO, len(data))
			for i, item := range data {
				mapped[i] = a.ResponseMapper.ToResponse(item)
			}
			response.Data = mapped
		}
	}

	return c.JSON(http.StatusOK, response)
}
