package action

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

var _ wool.PartiallyUpdate = (*UpdateAction[any, repository.Entity])(nil)

type UpdateAction[DTO any, Entity repository.Entity] struct {
	WriteRepository repository.WriteRepository[Entity]
	DTOFactory      DTOFactory[DTO]
	Mapper          Mapper[DTO, Entity]
}

func (a *UpdateAction[DTO, Entity]) PartiallyUpdate(ctx wool.Ctx) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	dto := a.DTOFactory.NewDTO()

	if err = ctx.Bind(dto); err != nil {
		return err
	}

	entity, err := a.Mapper.ToEntity(id, dto)
	if err != nil {
		return err
	}

	if err = a.WriteRepository.Save(ctx.Req().Context(), entity); err != nil {
		return err
	}

	return ctx.NoContent()
}
