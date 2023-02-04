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
	SourceID   uuid.UUID `json:"source_id,omitempty" bson:"source_id,omitempty"`
	Source     Source    `json:"source,omitempty" bson:"source,omitempty"`
	Lang       string    `json:"lang,omitempty" bson:"lang,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	ShortDesc  *string   `json:"short_desc,omitempty" bson:"short_desc,omitempty"`
	LongDesc   *string   `json:"long_desc,omitempty" bson:"long_desc,omitempty"`
	Guid       string    `json:"guid,omitempty" bson:"guid,omitempty"`
	Link       string    `json:"link,omitempty" bson:"link,omitempty"`
	Media      *[]Media  `json:"media,omitempty" bson:"media,omitempty"`
	PubDate    time.Time `json:"pub_date,omitempty" bson:"pub_date,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Authors    *[]string `json:"authors,omitempty" bson:"authors,omitempty"`
	Categories *[]string `json:"categories,omitempty" bson:"categories,omitempty"`
}

func (a *Article) EntityID() uuid.UUID {
	return a.ID
}

func (a *Article) Domain() string {
	return urlutil.SafeDomain(a.Link)
}

func (a *Article) SetShortDesc(shortDesc string) *Article {
	a.ShortDesc = &shortDesc
	return a
}

func (a *Article) SetLongDesc(longDesc string) *Article {
	a.LongDesc = &longDesc
	return a
}

func (a *Article) SetAuthors(authors []string) *Article {
	a.Authors = &authors
	return a
}

func (a *Article) SetCategories(categories []string) *Article {
	a.Categories = &categories
	return a
}

func (a *Article) SetMedia(media []Media) *Article {
	a.Media = &media
	return a
}

func (a *Article) FirstMedia(t MediaType) (m Media) {
	if a.Media != nil {
		for _, m = range *a.Media {
			if m.Type == t {
				return
			}
		}
	}
	return
}
