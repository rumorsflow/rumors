package pubsub

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"time"
)

type Article struct {
	ID         uuid.UUID     `json:"id,omitempty"`
	SourceID   uuid.UUID     `json:"source_id,omitempty"`
	Source     entity.Source `json:"source,omitempty"`
	Lang       string        `json:"lang,omitempty"`
	Title      string        `json:"title,omitempty"`
	ShortDesc  string        `json:"short_desc,omitempty"`
	Link       string        `json:"link,omitempty"`
	Image      string        `json:"url,omitempty"`
	PubDate    time.Time     `json:"pub_date,omitempty"`
	Categories []string      `json:"categories,omitempty"`
}

func FromEntity(e *entity.Article) Article {
	a := Article{
		ID:       e.ID,
		SourceID: e.SourceID,
		Source:   e.Source,
		Lang:     e.Lang,
		Title:    e.Title,
		Link:     e.Link,
		Image:    e.FirstMedia(entity.ImageType).URL,
		PubDate:  e.PubDate,
	}

	if e.ShortDesc != nil {
		a.ShortDesc = *e.ShortDesc
	}

	if e.Categories != nil {
		a.Categories = *e.Categories
	}

	return a
}
