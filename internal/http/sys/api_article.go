package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type MediaDTO struct {
	URL  string           `json:"url,omitempty" validate:"required,url"`
	Type entity.MediaType `json:"type,omitempty" validate:"required,max=10"`
	Meta map[string]any   `json:"meta,omitempty"`
}

func (dto MediaDTO) toEntity() entity.Media {
	return entity.Media{
		URL:  dto.URL,
		Type: dto.Type,
		Meta: dto.Meta,
	}
}

type UpdateArticleDTO struct {
	Title      string      `json:"title,omitempty" validate:"omitempty,max=254"`
	ShortDesc  *string     `json:"short_desc,omitempty"`
	LongDesc   *string     `json:"long_desc,omitempty"`
	Media      *[]MediaDTO `json:"media,omitempty" validate:"omitempty,dive"`
	Authors    *[]string   `json:"authors,omitempty" validate:"omitempty,dive,min=1,max=254"`
	Categories *[]string   `json:"categories,omitempty" validate:"omitempty,dive,min=1,max=254"`
}

func (dto UpdateArticleDTO) toEntity(id uuid.UUID) *entity.Article {
	a := &entity.Article{
		ID:         id,
		Title:      dto.Title,
		ShortDesc:  dto.ShortDesc,
		LongDesc:   dto.LongDesc,
		Authors:    dto.Authors,
		Categories: dto.Categories,
	}

	if dto.Media != nil {
		media := make([]entity.Media, len(*dto.Media))
		for i, m := range *dto.Media {
			media[i] = m.toEntity()
		}
		a.SetMedia(media)
	}

	return a
}

type ArticleActions struct {
	*action.ListAction[*entity.Article, any]
	*action.TakeAction[*entity.Article, any]
	*action.UpdateAction[*UpdateArticleDTO, *entity.Article]
	*action.DeleteAction[*entity.Article]
}

func NewArticleActions(
	read repository.ReadRepository[*entity.Article],
	write repository.WriteRepository[*entity.Article],
) *ArticleActions {
	return &ArticleActions{
		ListAction: &action.ListAction[*entity.Article, any]{ReadRepository: read},
		TakeAction: &action.TakeAction[*entity.Article, any]{ReadRepository: read},
		UpdateAction: &action.UpdateAction[*UpdateArticleDTO, *entity.Article]{
			WriteRepository: write,
			DTOFactory:      action.NewDTOFactory[*UpdateArticleDTO](),
			Mapper: action.RequestMapperFunc[*UpdateArticleDTO, *entity.Article](
				func(id uuid.UUID, dto *UpdateArticleDTO) (*entity.Article, error) {
					return dto.toEntity(id), nil
				},
			),
		},
		DeleteAction: &action.DeleteAction[*entity.Article]{WriteRepository: write},
	}
}
