package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"net/http"
)

var _ wool.Take = (*TakeAction[repository.Entity, any])(nil)

type TakeAction[Entity repository.Entity, DTO any] struct {
	ReadRepository repository.ReadRepository[Entity]
	ResponseMapper ResponseMapper[Entity, DTO]
}

func (a *TakeAction[Entity, DTO]) Take(c wool.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	entity, err := a.ReadRepository.FindByID(c.Req().Context(), id)
	if err != nil {
		return err
	}

	if a.ResponseMapper == nil {
		return c.JSON(http.StatusOK, entity)
	}

	return c.JSON(http.StatusOK, a.ResponseMapper.ToResponse(entity))
}
