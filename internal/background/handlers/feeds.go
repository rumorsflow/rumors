package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/daos"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/rs/zerolog"
	"strings"
)

type FeedsHandler struct {
	Notification notifications.Notification
	Dao          *daos.Dao
	Client       *asynq.Client
	Log          *zerolog.Logger
}

func (h *FeedsHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message

	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		h.Log.Error().Err(err).Msg("")
		return nil
	}

	switch task.Type() {
	case "feeds:crud":
		var name string

		switch strings.ToLower(Args(message.CommandArguments())[0]) {
		case "view", "show":
			name = "feeds:view"
		case "update", "edit":
			name = "feeds:update"
		default:
			name = "feeds:list"
		}

		_, err := h.Client.Enqueue(asynq.NewTask(name, task.Payload()))
		return err
	case "feeds:list":
		return h.list(ctx, message)
	case "feeds:view":
		return h.view(ctx, message)
	case "feeds:update":
		return h.update(ctx, message)
	}

	return nil
}

func (h *FeedsHandler) list(ctx context.Context, message tgbotapi.Message) error {
	i, s, f := Pagination(message.CommandArguments())
	var t *string
	if len(f) > 0 {
		t = &f[0]
	}

	data, err := h.Dao.FindFeeds(ctx, daos.FilterFeeds{Host: t}, i, s)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	if len(data) == 0 {
		h.Notification.Error(nil, "<b>No Feeds...</b>")
		return nil
	}

	var b strings.Builder
	for j, item := range data {
		b.WriteString(item.Line())

		if (j + 1) < len(data) {
			b.WriteString("\n")
		}
	}

	h.Notification.Send(nil, b.String())
	return nil
}

func (h *FeedsHandler) view(ctx context.Context, message tgbotapi.Message) error {
	id, _ := Id(message.CommandArguments())
	if id == 0 {
		h.Notification.Error(nil, "ID is required")
		return nil
	}

	feed, err := h.Dao.FindFeedById(ctx, id)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	h.Notification.Send(nil, feed.Info())
	return nil
}

func (h *FeedsHandler) update(ctx context.Context, message tgbotapi.Message) error {
	id, _ := Id(message.CommandArguments())
	if id == 0 {
		h.Notification.Error(nil, "ID is required")
		return nil
	}

	feed, err := h.Dao.FindFeedById(ctx, id)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	feed.Enabled = !feed.Enabled
	if err = h.Dao.Update(ctx, feed); err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	h.Notification.Send(nil, feed.Info())
	return nil
}
