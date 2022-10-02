package models

import "time"

type Room struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId    int64     `json:"chat_id" bson:"chat_id,omitempty"`
	Type      string    `json:"type" bson:"type,omitempty"`
	Title     string    `json:"title" bson:"title,omitempty"`
	Broadcast bool      `json:"broadcast" bson:"broadcast"`
	Deleted   bool      `json:"deleted" bson:"deleted"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
