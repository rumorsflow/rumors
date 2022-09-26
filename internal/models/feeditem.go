package models

import (
	"fmt"
	"github.com/iagapie/rumors/pkg/litedb/types"
	"strings"
)

type FeedItem struct {
	Id         int64           `db:"id" json:"id,omitempty"`
	FeedId     int64           `db:"feedId" json:"feedId"`
	Title      string          `db:"title" json:"title"`
	Desc       string          `db:"desc" json:"desc,omitempty"`
	Link       string          `db:"link" json:"link"`
	Guid       string          `db:"guid" json:"guid"`
	PubDate    types.DateTime  `db:"pubDate" json:"pubDate"`
	Created    types.DateTime  `db:"created" json:"created,omitempty"`
	Authors    types.JsonArray `db:"authors" json:"authors"`
	Categories types.JsonArray `db:"categories" json:"categories"`
}

func (*FeedItem) TableName() string {
	return "data.feed_items"
}

func (i *FeedItem) Table() string {
	return strings.Split(i.TableName(), ".")[1]
}

func (i *FeedItem) Line() string {
	return fmt.Sprintf("%s - <a href=\"%s\">link</a>", i.Title, i.Link)
}

func (i *FeedItem) Info() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>", i.Link, i.Title))

	for n, cat := range i.Categories {
		if n == 0 {
			b.WriteString("\n")
		}
		b.WriteString(cat.(string))
		if (n + 1) < len(i.Categories) {
			b.WriteString(", ")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(i.Desc)

	for n, cat := range i.Authors {
		if n == 0 {
			b.WriteString("\n")
		}
		b.WriteString(cat.(string))
		if (n + 1) < len(i.Authors) {
			b.WriteString(", ")
		}
	}

	return b.String()
}
