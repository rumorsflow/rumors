package models

import (
	"fmt"
	"github.com/iagapie/rumors/pkg/litedb/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Room struct {
	Id        int64          `db:"pk,id" json:"id"`
	Type      string         `db:"type" json:"type"`
	Title     string         `db:"title" json:"title"`
	Broadcast bool           `db:"broadcast" json:"broadcast"`
	Deleted   bool           `db:"deleted" json:"deleted"`
	Created   types.DateTime `db:"created" json:"created,omitempty"`
}

func (*Room) TableName() string {
	return "main.rooms"
}

func (r *Room) Line() string {
	return fmt.Sprintf("%d %s T[%s] B[%t] D[%t]", r.Id, r.Title, r.Type, r.Broadcast, r.Deleted)
}

func (r *Room) Info() string {
	text := fmt.Sprintf("<b>Room: %d</b>\n", r.Id)
	text += fmt.Sprintf("%s - %s\n\n", cases.Title(language.Und, cases.NoLower).String(r.Type), r.Title)
	text += fmt.Sprintf("Broadcast: %t\n", r.Broadcast)
	text += fmt.Sprintf("Deleted: %t\n", r.Deleted)
	text += fmt.Sprintf("Created: %s", r.Created.String())

	return text
}
