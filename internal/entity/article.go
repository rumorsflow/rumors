package entity

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/pkg/urlutil"
	"time"
)

type (
	Source    string
	MediaType string
)

const (
	ImageType MediaType = "image"
	VideoType MediaType = "video"
	AudioType MediaType = "audio"
)

const (
	FeedSource Source = "feed"
)

type Media struct {
	URL  string         `json:"url,omitempty" bson:"url,omitempty"`
	Type MediaType      `json:"type,omitempty" bson:"type,omitempty"`
	Meta map[string]any `json:"meta,omitempty" bson:"meta,omitempty"`
}

type Article struct {
	ID         uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Link       string    `json:"link,omitempty" bson:"link,omitempty"`
	SiteID     uuid.UUID `json:"site_id,omitempty" bson:"site_id,omitempty"`
	SourceID   uuid.UUID `json:"source_id,omitempty" bson:"source_id,omitempty"`
	Source     Source    `json:"source,omitempty" bson:"source,omitempty"`
	Lang       string    `json:"lang,omitempty" bson:"lang,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	ShortDesc  *string   `json:"short_desc,omitempty" bson:"short_desc,omitempty"`
	LongDesc   *string   `json:"long_desc,omitempty" bson:"long_desc,omitempty"`
	Media      *[]Media  `json:"media,omitempty" bson:"media,omitempty"`
	PubDate    time.Time `json:"pub_date,omitempty" bson:"pub_date,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Authors    *[]string `json:"authors,omitempty" bson:"authors,omitempty"`
	Categories *[]string `json:"categories,omitempty" bson:"categories,omitempty"`
}

func (e *Article) EntityID() uuid.UUID {
	return e.ID
}

func (e *Article) Domain() string {
	return urlutil.SafeDomain(e.Link)
}

func (e *Article) SetShortDesc(shortDesc string) *Article {
	e.ShortDesc = &shortDesc
	return e
}

func (e *Article) SetLongDesc(longDesc string) *Article {
	e.LongDesc = &longDesc
	return e
}

func (e *Article) SetAuthors(authors []string) *Article {
	e.Authors = &authors
	return e
}

func (e *Article) SetCategories(categories []string) *Article {
	e.Categories = &categories
	return e
}

func (e *Article) SetMedia(media []Media) *Article {
	e.Media = &media
	return e
}

func (e *Article) FirstMedia(t MediaType) (m Media) {
	if e.Media != nil {
		for _, m = range *e.Media {
			if m.Type == t {
				return
			}
		}
	}
	return
}
