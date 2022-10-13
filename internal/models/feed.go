package models

import (
	"github.com/rumorsflow/rumors/internal/pkg/url"
	"time"
)

type Feed struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	By        int64     `json:"by,omitempty" bson:"by,omitempty"`
	Languages []string  `json:"languages,omitempty" bson:"languages,omitempty"`
	Host      string    `json:"host,omitempty" bson:"host,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Link      string    `json:"link,omitempty" bson:"link,omitempty"`
	Enabled   *bool     `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (r *Feed) SetLink(link string) *Feed {
	r.Link = link
	r.Host = url.SafeDomain(link)
	return r
}

func (r *Feed) SetEnabled(enabled bool) *Feed {
	r.Enabled = &enabled
	return r
}
