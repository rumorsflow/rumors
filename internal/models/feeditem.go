package models

import (
	"github.com/rumorsflow/rumors/internal/pkg/url"
	"time"
)

type FeedItem struct {
	Id         string    `json:"id,omitempty" bson:"_id,omitempty"`
	FeedId     string    `json:"feed_id,omitempty" bson:"feed_id,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	Desc       *string   `json:"desc,omitempty" bson:"desc,omitempty"`
	Link       string    `json:"link,omitempty" bson:"link,omitempty"`
	Guid       string    `json:"guid,omitempty" bson:"guid,omitempty"`
	PubDate    time.Time `json:"pub_date,omitempty" bson:"pub_date,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Authors    []string  `json:"authors,omitempty" bson:"authors,omitempty"`
	Categories []string  `json:"categories,omitempty" bson:"categories,omitempty"`
}

func (i *FeedItem) Domain() string {
	return url.SafeDomain(i.Link)
}

func (i *FeedItem) SetDesc(desc string) *FeedItem {
	i.Desc = &desc
	return i
}
