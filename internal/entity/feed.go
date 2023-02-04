package entity

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/pkg/urlutil"
	"time"
)

type Feed struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Languages []string  `json:"languages,omitempty" bson:"languages,omitempty"`
	Host      string    `json:"host,omitempty" bson:"host,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Link      string    `json:"link,omitempty" bson:"link,omitempty"`
	Enabled   *bool     `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (f *Feed) EntityID() uuid.UUID {
	return f.ID
}

func (f *Feed) SetLink(link string) *Feed {
	f.Link = link
	if link != "" {
		f.Host = urlutil.SafeDomain(link)
	}
	return f
}

func (f *Feed) SetEnabled(enabled bool) *Feed {
	f.Enabled = &enabled
	return f
}
