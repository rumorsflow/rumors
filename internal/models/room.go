package models

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"time"
)

type Room struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId    int64     `json:"chat_id" bson:"chat_id,omitempty"`
	Type      string    `json:"type" bson:"type,omitempty"`
	Title     string    `json:"title" bson:"title,omitempty"`
	Broadcast bool      `json:"broadcast" bson:"broadcast"`
	Deleted   bool      `json:"deleted" bson:"deleted"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func (r *Room) Line() string {
	return fmt.Sprintf("%d %s T[%s] B[%t] D[%t]", r.ChatId, r.Title, r.Type, r.Broadcast, r.Deleted)
}

func (r *Room) Info() string {
	text := fmt.Sprintf("<b>Room: %d</b>\n", r.ChatId)
	text += fmt.Sprintf("%s - %s\n\n", cases.Title(language.Und, cases.NoLower).String(r.Type), r.Title)
	text += fmt.Sprintf("Broadcast: %t\n", r.Broadcast)
	text += fmt.Sprintf("Deleted: %t\n", r.Deleted)
	text += fmt.Sprintf("Created at: %s", r.CreatedAt.String())

	return text
}
