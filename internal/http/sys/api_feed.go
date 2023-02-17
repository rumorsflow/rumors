package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type CreateFeedDTO struct {
	SiteID  string `json:"site_id,omitempty" validate:"required,uuid4"`
	Link    string `json:"link,omitempty" validate:"required,url"`
	Enabled bool   `json:"enabled,omitempty"`
}

func (dto CreateFeedDTO) toEntity(id uuid.UUID) *entity.Feed {
	return (&entity.Feed{
		ID:     id,
		SiteID: uuid.MustParse(dto.SiteID),
		Link:   dto.Link,
	}).SetEnabled(dto.Enabled)
}

type UpdateFeedDTO struct {
	SiteID  string `json:"site_id,omitempty" validate:"omitempty,uuid4"`
	Link    string `json:"link,omitempty" validate:"omitempty,url"`
	Enabled *bool  `json:"enabled,omitempty"`
}

func (dto UpdateFeedDTO) toEntity(id uuid.UUID) *entity.Feed {
	siteID, _ := uuid.Parse(dto.SiteID)

	return &entity.Feed{
		ID:      id,
		SiteID:  siteID,
		Link:    dto.Link,
		Enabled: dto.Enabled,
	}
}

func NewFeedCRUD(
	read repository.ReadRepository[*entity.Feed],
	write repository.WriteRepository[*entity.Feed],
) action.CRUD {
	return action.NewCRUD[*CreateFeedDTO, *UpdateFeedDTO, *entity.Feed, any](
		read,
		write,
		action.NewDTOFactory[*CreateFeedDTO](),
		action.NewDTOFactory[*UpdateFeedDTO](),
		action.RequestMapperFunc[*CreateFeedDTO, *entity.Feed](func(id uuid.UUID, dto *CreateFeedDTO) (*entity.Feed, error) {
			return dto.toEntity(id), nil
		}),
		action.RequestMapperFunc[*UpdateFeedDTO, *entity.Feed](func(id uuid.UUID, dto *UpdateFeedDTO) (*entity.Feed, error) {
			return dto.toEntity(id), nil
		}),
		nil,
	)
}
