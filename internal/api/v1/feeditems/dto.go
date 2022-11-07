package feeditems

import "github.com/rumorsflow/rumors/internal/models"

type UpdateRequest struct {
	Title      string   `json:"title,omitempty" validate:"omitempty,max=254"`
	Desc       *string  `json:"desc,omitempty"`
	Authors    []string `json:"authors,omitempty" validate:"omitempty,dive,min=1,max=254"`
	Categories []string `json:"categories,omitempty" validate:"omitempty,dive,min=1,max=254"`
}

func (r UpdateRequest) FeedItem(id string) models.FeedItem {
	return models.FeedItem{
		Id:         id,
		Title:      r.Title,
		Desc:       r.Desc,
		Authors:    r.Authors,
		Categories: r.Categories,
	}
}
