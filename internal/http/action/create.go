package action

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
)

var _ wool.Create = (*CreateAction[any, repository.Entity])(nil)

type CreateAction[DTO any, Entity repository.Entity] struct {
	WriteRepository repository.WriteRepository[Entity]
	DTOFactory      DTOFactory[DTO]
	Mapper          RequestMapper[DTO, Entity]
}

func (a *CreateAction[DTO, Entity]) Create(c wool.Ctx) error {
	dto := a.DTOFactory.NewDTO()

	if err := c.Bind(dto); err != nil {
		return err
	}

	entity, err := a.Mapper.ToEntity(uuid.New(), dto)
	if err != nil {
		return err
	}

	if err = a.WriteRepository.Save(c.Req().Context(), entity); err != nil {
		return err
	}

	location := fmt.Sprintf(
		"%s://%s%s/%s",
		c.Req().URL.Scheme,
		c.Req().Host,
		c.Req().URL.Path,
		entity.EntityID().String(),
	)

	return c.Created(location)
}
