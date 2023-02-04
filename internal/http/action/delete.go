package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

var _ wool.Delete = (*DeleteAction[repository.Entity])(nil)

type DeleteAction[Entity repository.Entity] struct {
	WriteRepository repository.WriteRepository[Entity]
}

func (a *DeleteAction[Entity]) Delete(ctx wool.Ctx) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	if err = a.WriteRepository.Remove(ctx.Req().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent()
}
