package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"net/http"
)

var _ wool.Take = (*TakeAction[repository.Entity])(nil)

type TakeAction[Entity repository.Entity] struct {
	ReadRepository repository.ReadRepository[Entity]
}

func (a *TakeAction[Entity]) Take(ctx wool.Ctx) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	entity, err := a.ReadRepository.FindByID(ctx.Req().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, entity)
}
