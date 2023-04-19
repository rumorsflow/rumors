package entity

import (
	"github.com/google/uuid"
	"time"
)

type Site struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Domain    string    `json:"domain,omitempty" bson:"domain,omitempty"`
	Favicon   string    `json:"favicon,omitempty" bson:"favicon,omitempty"`
	Languages []string  `json:"languages,omitempty" bson:"languages,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Enabled   *bool     `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (e *Site) EntityID() uuid.UUID {
	return e.ID
}

func (e *Site) SetEnabled(enabled bool) *Site {
	e.Enabled = &enabled
	return e
}
