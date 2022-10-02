package models

import "time"

type Feed struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	By        int64     `json:"by" bson:"by,omitempty"`
	Lang      string    `json:"lang,omitempty" bson:"lang,omitempty" validate:"omitempty,bcp47_language_tag"`
	Host      string    `json:"host,omitempty" bson:"host,omitempty" validate:"omitempty,fqdn"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Link      string    `json:"link" bson:"link,omitempty" validate:"required,url"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
