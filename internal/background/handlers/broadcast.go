package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog"
	"strings"
	"sync"
)

type BroadcastHandler struct {
	Notification notifications.Notification
	FeedStorage  storage.FeedStorage
	RoomStorage  storage.RoomStorage
	Log          *zerolog.Logger
}

func (h *BroadcastHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var items []models.FeedItem
	if err := json.Unmarshal(task.Payload(), &items); err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	feed, err := h.FeedStorage.FindById(ctx, items[0].FeedId)
	if err != nil {
		h.Notification.Err(nil, err)
		return nil
	}

	var b strings.Builder
	b.WriteString("<b>")
	b.WriteString(feed.Host)
	b.WriteString("</b>\n\n")

	for j, item := range items {
		if len(items) > 1 {
			b.WriteString(fmt.Sprintf("<b>%d.</b> ", j+1))
		}
		if len(items) > 3 {
			b.WriteString(item.Line())
		} else {
			b.WriteString(item.Info())
		}
		b.WriteString("\n\n")
	}

	text := strings.TrimSuffix(b.String(), "\n\n")

	var wg sync.WaitGroup

	broadcast := true
	deleted := false
	filter := storage.FilterRooms{Broadcast: &broadcast, Deleted: &deleted}
	var index uint64 = 0
	for ; ; index += 20 {
		rooms, err := h.RoomStorage.Find(ctx, filter, index, 20)
		if err != nil {
			h.Notification.Err(nil, err)
			return nil
		}
		wg.Add(len(rooms))
		for _, room := range rooms {
			go func(room models.Room) {
				defer wg.Done()
				h.Notification.Send(room.ChatId, text)
			}(room)
		}
		if len(rooms) < 20 {
			break
		}
	}

	wg.Wait()

	return nil
}
