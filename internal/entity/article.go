package entity

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/pkg/util"
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
	FeedSource    Source = "feed"
	SitemapSource Source = "sitemap"

	ArticleCollection = "articles"
)

type Media struct {
	URL  string         `json:"url,omitempty" bson:"url,omitempty"`
	Type MediaType      `json:"type,omitempty" bson:"type,omitempty"`
	Meta map[string]any `json:"meta,omitempty" bson:"meta,omitempty"`
}

type Article struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Link      string    `json:"link,omitempty" bson:"link,omitempty"`
	SiteID    uuid.UUID `json:"site_id,omitempty" bson:"site_id,omitempty"`
	Source    Source    `json:"source,omitempty" bson:"source,omitempty"`
	Lang      string    `json:"lang,omitempty" bson:"lang,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Desc      *string   `json:"desc,omitempty" bson:"short_desc,omitempty"`
	Media     *[]Media  `json:"media,omitempty" bson:"media,omitempty"`
	PubDate   time.Time `json:"pub_date,omitempty" bson:"pub_date,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (e *Article) Tags() []string {
	return []string{ArticleCollection, e.ID.String()}
}

func (e *Article) EntityID() uuid.UUID {
	return e.ID
}

func (e *Article) Domain() string {
	return util.SafeDomain(e.Link)
}

func (e *Article) SetDesc(desc string) *Article {
	e.Desc = &desc
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
