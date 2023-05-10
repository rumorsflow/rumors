package model

import (
	"github.com/google/uuid"
	"github.com/mergestat/timediff"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"time"
)

type Article struct {
	ID         uuid.UUID `json:"id,omitempty"`
	SiteID     uuid.UUID `json:"site_id,omitempty"`
	Lang       string    `json:"lang,omitempty"`
	Title      string    `json:"title,omitempty"`
	ShortDesc  string    `json:"short_desc,omitempty"`
	Link       string    `json:"link,omitempty"`
	Image      string    `json:"image,omitempty"`
	PubDate    time.Time `json:"pub_date,omitempty"`
	PubDiff    string    `json:"pub_diff,omitempty"`
	Categories []string  `json:"categories,omitempty"`
}

func ArticleFromEntity(e *entity.Article) Article {
	a := Article{
		ID:      e.ID,
		SiteID:  e.SiteID,
		Lang:    e.Lang,
		Title:   e.Title,
		Link:    e.Link,
		Image:   e.FirstMedia(entity.ImageType).URL,
		PubDate: e.CreatedAt,
		PubDiff: timediff.TimeDiff(e.CreatedAt, timediff.WithStartTime(time.Now().UTC())),
	}

	if e.ShortDesc != nil {
		a.ShortDesc = *e.ShortDesc
	}

	if e.Categories != nil {
		a.Categories = *e.Categories
	}

	return a
}
