package entity

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/pkg/urlutil"
	"time"
)

type Feed struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	SiteID    uuid.UUID `json:"site_id,omitempty" bson:"site_id,omitempty"`
	Link      string    `json:"link,omitempty" bson:"link,omitempty"`
	Enabled   *bool     `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (e *Feed) EntityID() uuid.UUID {
	return e.ID
}

func (e *Feed) Domain() string {
	return urlutil.SafeDomain(e.Link)
}

func (e *Feed) SetEnabled(enabled bool) *Feed {
	e.Enabled = &enabled
	return e
}
