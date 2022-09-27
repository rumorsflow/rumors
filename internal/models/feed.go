package models

import (
	"fmt"
	"time"
)

type Feed struct {
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	By        int64     `json:"by" bson:"by,omitempty"`
	Host      string    `json:"host" bson:"host,omitempty"`
	Title     string    `json:"title" bson:"title,omitempty"`
	Link      string    `json:"link" bson:"link,omitempty"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func (f *Feed) Line() string {
	return fmt.Sprintf("%s - %s E[%t]", f.Id, f.Link, f.Enabled)
}

func (f *Feed) Info() string {
	text := fmt.Sprintf("<b>Feed: %s</b>\n\n", f.Id)
	text += fmt.Sprintf("<b>%s</b> - <u><a href=\"%s\">%s</a></u>\n\n", f.Host, f.Link, f.Title)
	text += fmt.Sprintf("By: %d\n", f.By)
	text += fmt.Sprintf("Enabled: %t\n", f.Enabled)
	text += fmt.Sprintf("Created at: %s", f.CreatedAt.String())

	return text
}
