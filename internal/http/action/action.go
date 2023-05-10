package action

import (
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
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

type RequestMapper[DTO any, Entity repository.Entity] interface {
	ToEntity(id uuid.UUID, dto DTO) (Entity, error)
}

type RequestMapperFunc[DTO any, Entity repository.Entity] func(id uuid.UUID, dto DTO) (Entity, error)

func (m RequestMapperFunc[DTO, Entity]) ToEntity(id uuid.UUID, dto DTO) (Entity, error) {
	return m(id, dto)
}

type ResponseMapper[Entity repository.Entity, DTO any] interface {
	ToResponse(entity Entity) DTO
}

type ResponseMapperFunc[Entity repository.Entity, DTO any] func(entity Entity) DTO

func (m ResponseMapperFunc[Entity, DTO]) ToResponse(entity Entity) DTO {
	return m(entity)
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

type crud[CreateDTO any, UpdateDTO any, Entity repository.Entity, ResponseDTO any] struct {
	*ListAction[Entity, ResponseDTO]
	*TakeAction[Entity, ResponseDTO]
	*CreateAction[CreateDTO, Entity]
	*UpdateAction[UpdateDTO, Entity]
	*DeleteAction[Entity]
}

func NewCRUD[CreateDTO any, UpdateDTO any, Entity repository.Entity, ResponseDTO any](
	read repository.ReadRepository[Entity],
	write repository.WriteRepository[Entity],
	createDTO DTOFactory[CreateDTO],
	updateDTO DTOFactory[UpdateDTO],
	createMapper RequestMapper[CreateDTO, Entity],
	updateMapper RequestMapper[UpdateDTO, Entity],
	responseMapper ResponseMapper[Entity, ResponseDTO],
) CRUD {
	return &crud[CreateDTO, UpdateDTO, Entity, ResponseDTO]{
		ListAction:   &ListAction[Entity, ResponseDTO]{ReadRepository: read, ResponseMapper: responseMapper},
		TakeAction:   &TakeAction[Entity, ResponseDTO]{ReadRepository: read, ResponseMapper: responseMapper},
		CreateAction: &CreateAction[CreateDTO, Entity]{WriteRepository: write, DTOFactory: createDTO, Mapper: createMapper},
		UpdateAction: &UpdateAction[UpdateDTO, Entity]{WriteRepository: write, DTOFactory: updateDTO, Mapper: updateMapper},
		DeleteAction: &DeleteAction[Entity]{WriteRepository: write},
	}
}
