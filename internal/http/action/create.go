package action

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

var _ wool.Create = (*CreateAction[any, repository.Entity])(nil)

type CreateAction[DTO any, Entity repository.Entity] struct {
	WriteRepository repository.WriteRepository[Entity]
	DTOFactory      DTOFactory[DTO]
	Mapper          Mapper[DTO, Entity]
}

func (a *CreateAction[DTO, Entity]) Create(ctx wool.Ctx) error {
	dto := a.DTOFactory.NewDTO()

	if err := ctx.Bind(dto); err != nil {
		return err
	}

	entity, err := a.Mapper.ToEntity(uuid.New(), dto)
	if err != nil {
		return err
	}

	if err = a.WriteRepository.Save(ctx.Req().Context(), entity); err != nil {
		return err
	}

	location := fmt.Sprintf(
		"%s://%s%s/%s",
		ctx.Req().URL.Scheme,
		ctx.Req().Host,
		ctx.Req().URL.Path,
		entity.EntityID().String(),
	)

	return ctx.Created(location)
}
