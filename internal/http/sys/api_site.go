package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
)

type CreateSiteDTO struct {
	Domain    string   `json:"domain,omitempty" validate:"required,fqdn"`
	Favicon   string   `json:"favicon,omitempty" validate:"required,url"`
	Languages []string `json:"languages,omitempty" validate:"required,min=1,dive,bcp47_language_tag"`
	Title     string   `json:"title,omitempty" validate:"required,max=254"`
	Enabled   bool     `json:"enabled,omitempty"`
}

func (dto CreateSiteDTO) toEntity(id uuid.UUID) *entity.Site {
	return (&entity.Site{
		ID:        id,
		Domain:    dto.Domain,
		Favicon:   dto.Favicon,
		Languages: dto.Languages,
		Title:     dto.Title,
	}).SetEnabled(dto.Enabled)
}

type UpdateSiteDTO struct {
	Domain    string   `json:"domain,omitempty" validate:"omitempty,fqdn"`
	Favicon   string   `json:"favicon,omitempty" validate:"required,url"`
	Languages []string `json:"languages,omitempty" validate:"omitempty,dive,bcp47_language_tag"`
	Title     string   `json:"title,omitempty" validate:"omitempty,max=254"`
	Enabled   *bool    `json:"enabled,omitempty"`
}

func (dto UpdateSiteDTO) toEntity(id uuid.UUID) *entity.Site {
	return &entity.Site{
		ID:        id,
		Domain:    dto.Domain,
		Favicon:   dto.Favicon,
		Languages: dto.Languages,
		Title:     dto.Title,
		Enabled:   dto.Enabled,
	}
}

func NewSiteCRUD(read repository.ReadRepository[*entity.Site], write repository.WriteRepository[*entity.Site]) action.CRUD {
	return action.NewCRUD[*CreateSiteDTO, *UpdateSiteDTO, *entity.Site, any](
		read,
		write,
		action.NewDTOFactory[*CreateSiteDTO](),
		action.NewDTOFactory[*UpdateSiteDTO](),
		action.RequestMapperFunc[*CreateSiteDTO, *entity.Site](func(id uuid.UUID, dto *CreateSiteDTO) (*entity.Site, error) {
			return dto.toEntity(id), nil
		}),
		action.RequestMapperFunc[*UpdateSiteDTO, *entity.Site](func(id uuid.UUID, dto *UpdateSiteDTO) (*entity.Site, error) {
			return dto.toEntity(id), nil
		}),
		nil,
	)
}
