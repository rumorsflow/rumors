package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/daos"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/rs/zerolog"
	"strings"
)

type RumorsHandler struct {
	Notification notifications.Notification
	Dao          *daos.Dao
	Log          *zerolog.Logger
}

func (h *RumorsHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message

	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	return h.list(ctx, message)
}

func (h *RumorsHandler) list(ctx context.Context, message tgbotapi.Message) error {
	i, s, f := Pagination(message.CommandArguments())

	feeds := make(map[int64]models.Feed)

	var feedIds []int64

	if len(f) > 0 {
		if data, err := h.Dao.FindFeeds(ctx, daos.FilterFeeds{Host: &f[0]}, 0, 50); err == nil {
			for _, item := range data {
				feedIds = append(feedIds, item.Id)
				feeds[item.Id] = item
			}
		}
	}

	data, err := h.Dao.FindFeedItems(ctx, daos.FilterFeedItems{FeedIds: feedIds}, i, s)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	if len(data) == 0 {
		h.Notification.Send(message.Chat.ID, "<b>No Rumors...</b>")
		return nil
	}

	group := make(map[string][]models.FeedItem)

	for _, item := range data {
		if _, ok := feeds[item.FeedId]; !ok {
			if feed, err := h.Dao.FindFeedById(ctx, item.FeedId); err == nil {
				feeds[feed.Id] = *feed
			} else {
				continue
			}
		}

		feed := feeds[item.FeedId]

		if _, ok := group[feed.Host]; ok {
			group[feed.Host] = append(group[feed.Host], item)
		} else {
			group[feed.Host] = []models.FeedItem{item}
		}
	}

	var b strings.Builder
	for host, items := range group {
		b.WriteString("<b>")
		b.WriteString(host)
		b.WriteString("</b>\n\n")

		if len(items) > 3 {
			for j, item := range items {
				b.WriteString(fmt.Sprintf("<b>%d.</b> ", j+1))
				b.WriteString(item.Line())
				b.WriteString("\n\n")
			}
		} else {
			for _, item := range items {
				b.WriteString(item.Info())
				b.WriteString("\n\n")
			}
		}

		b.WriteString("\n\n")
	}

	text := strings.TrimSuffix(b.String(), "\n\n\n\n")

	h.Notification.Send(message.Chat.ID, text)

	return nil
}
