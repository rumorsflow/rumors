package models

import (
	"fmt"
	"github.com/iagapie/rumors/pkg/litedb/types"
)

type Feed struct {
	Id      int64          `db:"id" json:"id,omitempty"`
	By      int64          `db:"by" json:"by"`
	Host    string         `db:"host" json:"host"`
	Title   string         `db:"title" json:"title"`
	Link    string         `db:"link" json:"link"`
	Enabled bool           `db:"enabled" json:"enabled"`
	Created types.DateTime `db:"created" json:"created,omitempty"`
}

func (*Feed) TableName() string {
	return "main.feeds"
}

func (f *Feed) Line() string {
	return fmt.Sprintf("%d - %s E[%t]", f.Id, f.Link, f.Enabled)
}

func (f *Feed) Info() string {
	text := fmt.Sprintf("<b>Feed: %d</b>\n\n", f.Id)
	text += fmt.Sprintf("<b>%s</b> - <u><a href=\"%s\">%s</a></u>\n\n", f.Host, f.Link, f.Title)
	text += fmt.Sprintf("By: %d\n", f.By)
	text += fmt.Sprintf("Enabled: %t\n", f.Enabled)
	text += fmt.Sprintf("Created: %s", f.Created.String())

	return text
}
