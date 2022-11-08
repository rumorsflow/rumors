package parser

import (
	"context"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/pkg/str"
	"go.uber.org/zap"
	"strings"
	"time"
	"unicode/utf8"
)

var _ FeedParser = (*Plugin)(nil)

type FeedParser interface {
	Parse(ctx context.Context, feed models.Feed) ([]models.FeedItem, error)
}

func (p *Plugin) Parse(ctx context.Context, feed models.Feed) ([]models.FeedItem, error) {
	parsed, err := gofeed.NewParser().ParseURLWithContext(feed.Link, ctx)
	if err != nil {
		return nil, err
	}

	var data []models.FeedItem

	for _, item := range parsed.Items {
		link := item.Link
		if link == "" && len(item.Links) > 0 {
			link = item.Links[0]
		}

		guid := item.GUID

		if link == "" || guid == "" {
			p.log.Warn("link or guid is empty", zap.Any("parsed_item", item))
			continue
		}

		desc := strings.TrimSpace(item.Description)
		if desc == "" {
			desc = item.Content
		}
		desc = str.StripHTMLTags(desc)

		title := str.StripHTMLTags(item.Title)
		if title == "" {
			if desc == "" {
				p.log.Warn("title and desc are empty", zap.Any("parsed_item", item))
				continue
			}

			index := strings.Index(desc, ". ")
			if index != -1 && index > 0 {
				title = desc[:index]
			} else if utf8.RuneCountInString(desc) > 100 {
				title = strings.TrimSuffix(string([]rune(desc)[:97]), ".") + "..."
			} else {
				title = desc
			}
		}

		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubDate = *item.UpdatedParsed
		}

		authors := make([]string, 0)
		for _, author := range item.Authors {
			if author != nil {
				if name := str.StripHTMLTags(author.Name); name != "" {
					authors = append(authors, name)
				}
			}
		}

		categories := make([]string, len(item.Categories))
		for i, c := range item.Categories {
			if c = str.StripHTMLTags(c); c != "" {
				categories[i] = c
			}
		}

		var d *string
		if desc != "" {
			d = &desc
		}

		data = append(data, models.FeedItem{
			Id:         uuid.NewString(),
			FeedId:     feed.Id,
			Title:      item.Title,
			Desc:       d,
			Link:       link,
			Guid:       guid,
			PubDate:    pubDate,
			Authors:    authors,
			Categories: categories,
		})
	}

	return data, nil
}
