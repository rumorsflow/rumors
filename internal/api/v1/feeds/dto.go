package feeds

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/internal/models"
)

type CreateRequest struct {
	Languages []string `json:"languages,omitempty" validate:"required,min=1,dive,bcp47_language_tag"`
	Host      string   `json:"host,omitempty" validate:"omitempty,max=254"`
	Title     string   `json:"title,omitempty" validate:"required,max=254"`
	Link      string   `json:"link,omitempty" validate:"required,url"`
	Enabled   bool     `json:"enabled,omitempty"`
}

func (r CreateRequest) Feed(by int64) models.Feed {
	model := models.Feed{
		Id:        uuid.NewString(),
		By:        by,
		Languages: r.Languages,
		Title:     r.Title,
	}
	model.SetLink(r.Link).SetEnabled(r.Enabled)

	if r.Host != "" {
		model.Host = r.Host
	}

	return model
}

type UpdateRequest struct {
	Languages []string `json:"languages,omitempty" validate:"omitempty,dive,bcp47_language_tag"`
	Host      string   `json:"host,omitempty" validate:"omitempty,max=254"`
	Title     string   `json:"title,omitempty" validate:"omitempty,max=254"`
	Link      string   `json:"link,omitempty" validate:"omitempty,url"`
	Enabled   *bool    `json:"enabled,omitempty"`
}

func (r UpdateRequest) Feed(id string) models.Feed {
	model := models.Feed{
		Id:        id,
		Languages: r.Languages,
		Title:     r.Title,
		Enabled:   r.Enabled,
	}

	if r.Link != "" {
		model.SetLink(r.Link)
	}

	if r.Host != "" {
		model.Host = r.Host
	}

	return model
}
