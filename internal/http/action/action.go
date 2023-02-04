package action

import (
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"reflect"
)

type DTOFactory[DTO any] interface {
	NewDTO() DTO
}

type DTOFactoryFunc[DTO any] func() DTO

func (f DTOFactoryFunc[DTO]) NewDTO() DTO {
	return f()
}

func NewDTOFactory[DTO any]() DTOFactory[DTO] {
	return DTOFactoryFunc[DTO](func() (dto DTO) {
		d := reflect.ValueOf(dto)
		if d.Kind() == reflect.Ptr && d.IsNil() {
			dto = reflect.New(d.Type().Elem()).Interface().(DTO)
		}
		return
	})
}

type Mapper[DTO any, Entity repository.Entity] interface {
	ToEntity(id uuid.UUID, dto DTO) (Entity, error)
}

type MapperFunc[DTO any, Entity repository.Entity] func(id uuid.UUID, dto DTO) (Entity, error)

func (m MapperFunc[DTO, Entity]) ToEntity(id uuid.UUID, dto DTO) (Entity, error) {
	return m(id, dto)
}

func parseID(ctx wool.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Req().PathParamID())
	if err != nil {
		return uuid.Nil, wool.NewErrBadRequest(err, "id param is not valid")
	}
	return id, nil
}

type CRUD interface {
	wool.List
	wool.Take
	wool.Create
	wool.PartiallyUpdate
	wool.Delete
}

type crud[CreateDTO any, UpdateDTO any, Entity repository.Entity] struct {
	*ListAction[Entity]
	*TakeAction[Entity]
	*CreateAction[CreateDTO, Entity]
	*UpdateAction[UpdateDTO, Entity]
	*DeleteAction[Entity]
}

func NewCRUD[CreateDTO any, UpdateDTO any, Entity repository.Entity](
	read repository.ReadRepository[Entity],
	write repository.WriteRepository[Entity],
	createDTO DTOFactory[CreateDTO],
	updateDTO DTOFactory[UpdateDTO],
	createMapper Mapper[CreateDTO, Entity],
	updateMapper Mapper[UpdateDTO, Entity],
) CRUD {
	return &crud[CreateDTO, UpdateDTO, Entity]{
		ListAction:   &ListAction[Entity]{ReadRepository: read},
		TakeAction:   &TakeAction[Entity]{ReadRepository: read},
		CreateAction: &CreateAction[CreateDTO, Entity]{WriteRepository: write, DTOFactory: createDTO, Mapper: createMapper},
		UpdateAction: &UpdateAction[UpdateDTO, Entity]{WriteRepository: write, DTOFactory: updateDTO, Mapper: updateMapper},
		DeleteAction: &DeleteAction[Entity]{WriteRepository: write},
	}
}
