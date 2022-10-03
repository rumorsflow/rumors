package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/str"
	"github.com/iagapie/rumors/pkg/url"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type FeedImporterHandler struct {
	Client *asynq.Client
}

func (h *FeedImporterHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	l := log.Ctx(ctx)

	var feed models.Feed
	if err := json.Unmarshal(task.Payload(), &feed); err != nil {
		l.Error().Err(err).Msg("error due to unmarshal task payload")
		return nil
	}

	parsed, err := gofeed.NewParser().ParseURLWithContext(feed.Link, ctx)
	if err != nil {
		l.Error().Err(err).Msg("error due to parse feed")
		return nil
	}

	if feed.Id == "" {
		defer func() {
			recover()
		}()
		feed.Id = uuid.NewString()
		feed.Host = url.MustDomain(feed.Link)
		feed.Title = parsed.Title
		feed.Link = parsed.FeedLink
		feed.CreatedAt = time.Now().UTC()

		payload, err := json.Marshal(feed)
		if err != nil {
			l.Error().Err(err).Str("feedLink", feed.Link).Msg("error due to marshal feed")
			return nil
		}

		if _, err = h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskFeedSave, payload)); err != nil {
			l.Error().Err(err).RawJSON("feed", payload).Msg("error due to enqueue feed")
			return err
		}

		if !feed.Enabled {
			return nil
		}
	}

	for _, item := range parsed.Items {
		feedItem := toFeedItem(item, feed.Id)
		if feedItem == nil {
			continue
		}

		payload, err := json.Marshal(feedItem)
		if err != nil {
			l.Error().Err(err).Interface("feed", feed).Str("feedItemLink", feedItem.Link).Msg("error due to marshal feed item")
			continue
		}

		id := asynq.TaskID(feedItem.PubDate.String() + feedItem.Guid)

		if _, err = h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskFeedItemSave, payload), id); err != nil {
			if !errors.Is(err, asynq.ErrTaskIDConflict) {
				l.Error().Err(err).Interface("feed", feed).RawJSON("feedItem", payload).Msg("error due to enqueue feed item")
			}
		}
	}

	return nil
}

func toFeedItem(item *gofeed.Item, feedId string) *models.FeedItem {
	link := item.Link
	if link == "" && len(item.Links) > 0 {
		link = item.Links[0]
	}

	guid := item.GUID

	if link == "" || guid == "" {
		return nil
	}

	desc := strings.TrimSpace(item.Description)
	if desc == "" {
		desc = item.Content
	}
	desc = str.StripHTMLTags(desc)

	title := str.StripHTMLTags(item.Title)
	if title == "" {
		if desc == "" {
			return nil
		}

		index := strings.Index(desc, ".")
		if index != -1 && index > 0 {
			title = desc[0:index]
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
			authors = append(authors, str.StripHTMLTags(author.Name))
		}
	}

	categories := make([]string, len(item.Categories))
	for i, c := range item.Categories {
		categories[i] = str.StripHTMLTags(c)
	}

	return &models.FeedItem{
		Id:         uuid.NewString(),
		FeedId:     feedId,
		Title:      item.Title,
		Desc:       desc,
		Link:       link,
		Guid:       guid,
		PubDate:    pubDate,
		CreatedAt:  time.Now().UTC(),
		Authors:    authors,
		Categories: categories,
	}
}
