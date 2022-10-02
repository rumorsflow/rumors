package models

import (
	"github.com/iagapie/rumors/pkg/url"
	"time"
)

type FeedItem struct {
	Id         string    `json:"id,omitempty" bson:"_id,omitempty"`
	FeedId     string    `json:"feed_id" bson:"feed_id,omitempty"`
	Title      string    `json:"title" bson:"title,omitempty"`
	Desc       string    `json:"desc,omitempty" bson:"desc,omitempty"`
	Link       string    `json:"link" bson:"link,omitempty"`
	Guid       string    `json:"guid" bson:"guid,omitempty"`
	PubDate    time.Time `json:"pub_date" bson:"pub_date,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Authors    []string  `json:"authors" bson:"authors,omitempty"`
	Categories []string  `json:"categories" bson:"categories,omitempty"`
}

func (i *FeedItem) Domain() string {
	domain, _ := url.Domain(i.Link)
	return domain
}
