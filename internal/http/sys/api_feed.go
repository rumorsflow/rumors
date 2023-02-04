package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type CreateFeedDTO struct {
	Languages []string `json:"languages,omitempty" validate:"required,min=1,dive,bcp47_language_tag"`
	Title     string   `json:"title,omitempty" validate:"required,max=254"`
	Link      string   `json:"link,omitempty" validate:"required,url"`
	Enabled   bool     `json:"enabled,omitempty"`
}

func (dto CreateFeedDTO) toEntity(id uuid.UUID) *entity.Feed {
	return (&entity.Feed{
		ID:        id,
		Languages: dto.Languages,
		Title:     dto.Title,
	}).SetLink(dto.Link).SetEnabled(dto.Enabled)
}

type UpdateFeedDTO struct {
	Languages []string `json:"languages,omitempty" validate:"omitempty,dive,bcp47_language_tag"`
	Title     string   `json:"title,omitempty" validate:"omitempty,max=254"`
	Link      string   `json:"link,omitempty" validate:"omitempty,url"`
	Enabled   *bool    `json:"enabled,omitempty"`
}

func (dto UpdateFeedDTO) toEntity(id uuid.UUID) *entity.Feed {
	return (&entity.Feed{
		ID:        id,
		Languages: dto.Languages,
		Title:     dto.Title,
		Enabled:   dto.Enabled,
	}).SetLink(dto.Link)
}

func NewFeedCRUD(
	read repository.ReadRepository[*entity.Feed],
	write repository.WriteRepository[*entity.Feed],
) action.CRUD {
	return action.NewCRUD[*CreateFeedDTO, *UpdateFeedDTO, *entity.Feed](
		read,
		write,
		action.NewDTOFactory[*CreateFeedDTO](),
		action.NewDTOFactory[*UpdateFeedDTO](),
		action.MapperFunc[*CreateFeedDTO, *entity.Feed](func(id uuid.UUID, dto *CreateFeedDTO) (*entity.Feed, error) {
			return dto.toEntity(id), nil
		}),
		action.MapperFunc[*UpdateFeedDTO, *entity.Feed](func(id uuid.UUID, dto *UpdateFeedDTO) (*entity.Feed, error) {
			return dto.toEntity(id), nil
		}),
	)
}
