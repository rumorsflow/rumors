package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"strings"
	"sync"
	"time"
)

type FeedsParserHandler struct {
	Notification    notifications.Notification
	FeedStorage     storage.FeedStorage
	FeedItemStorage storage.FeedItemStorage
	Client          *asynq.Client
	Log             *zerolog.Logger
}

func (h *FeedsParserHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	var wg sync.WaitGroup

	enabled := true
	filter := storage.FilterFeeds{Enabled: &enabled}
	var index uint64 = 0
	for ; ; index += 20 {
		data, err := h.FeedStorage.Find(ctx, filter, index, 20)
		if err != nil {
			h.Notification.Err(nil, err)
			return nil
		}

		for _, item := range data {
			wg.Add(1)
			go func(feed models.Feed) {
				defer wg.Done()
				h.parse(feed)
			}(item)
		}

		if len(data) < 20 {
			break
		}
	}

	wg.Wait()

	return nil
}

func (h *FeedsParserHandler) parse(feed models.Feed) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parsed, err := gofeed.NewParser().ParseURLWithContext(feed.Link, ctx)
	if err != nil {
		h.Notification.Err(nil, err)
		h.Notification.Send(nil, feed.Info())
		return
	}

	group := asynq.Group(feed.Host)

	for _, item := range parsed.Items {
		link := item.Link
		if link == "" && len(item.Links) > 0 {
			link = item.Links[0]
		}

		guid := item.GUID

		if link == "" || guid == "" {
			continue
		}

		desc := item.Description
		if desc == "" {
			desc = item.Content
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
				authors = append(authors, strings.TrimSpace(author.Name))
			}
		}

		categories := make([]string, len(item.Categories))
		for i, c := range item.Categories {
			categories[i] = strings.TrimSpace(c)
		}

		fItem := models.FeedItem{
			Id:         uuid.NewString(),
			FeedId:     feed.Id,
			Title:      item.Title,
			Desc:       desc,
			Link:       link,
			Guid:       guid,
			PubDate:    pubDate,
			CreatedAt:  time.Now().UTC(),
			Authors:    authors,
			Categories: categories,
		}

		if err = h.FeedItemStorage.Save(ctx, fItem); err != nil {
			h.Log.Debug().Err(err).Interface("feed", feed).Interface("feed_item", fItem).Msg("")
			continue
		}

		if data, err := json.Marshal(fItem); err == nil {
			task := asynq.NewTask("aggregated:item", data)
			if _, err = h.Client.Enqueue(task, asynq.Queue("broadcast"), group); err != nil {
				h.Notification.Err(nil, err)
			}
		}
	}
}
