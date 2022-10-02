package handlers

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog/log"
)

type FeedItemAggregatedHandler struct {
	Storage storage.RoomStorage
	Client  *asynq.Client
}

func (h *FeedItemAggregatedHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	if len(task.Payload()) <= 2 {
		return nil
	}

	l := log.Ctx(ctx)

	var index uint64 = 0
	var size uint32 = 50
	filter := new(storage.FilterRooms).SetBroadcast(true).SetDeleted(false)

	for {
		rooms, err := h.Storage.Find(ctx, filter, index, size)
		if err != nil {
			l.Error().Err(err).Msg("error due to find rooms")
			return nil
		}

		for _, room := range rooms {
			payload := append(Int64ToBytes(room.ChatId), task.Payload()...)
			if _, err = h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskFeedItemBroadcast, payload)); err != nil {
				l.Error().Err(err).Msg("error due to enqueue aggregated payload")
				return err
			}
		}

		if len(rooms) < int(size) {
			break
		}

		index += uint64(size)
	}
	return nil
}
