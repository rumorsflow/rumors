package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
)

var _ wool.Delete = (*DeleteAction[repository.Entity])(nil)

type DeleteAction[Entity repository.Entity] struct {
	WriteRepository repository.WriteRepository[Entity]
}

func (a *DeleteAction[Entity]) Delete(c wool.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	if err = a.WriteRepository.Remove(c.Req().Context(), id); err != nil {
		return err
	}

	return c.NoContent()
}
