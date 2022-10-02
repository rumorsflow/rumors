package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/rs/zerolog/log"
)

type FeedSchedulerHandler struct {
	Storage storage.FeedStorage
	Client  *asynq.Client
}

func (h *FeedSchedulerHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	l := log.Ctx(ctx)

	var index uint64 = 0
	var size uint32 = 50
	filter := new(storage.FilterFeeds).SetEnabled(true)

	for {
		feeds, err := h.Storage.Find(ctx, filter, index, size)
		if err != nil {
			l.Error().Err(err).Uint64("index", index).Msg("error due to find feeds")
			return nil
		}

		for _, feed := range feeds {
			payload, err := json.Marshal(feed)
			if err != nil {
				l.Error().Err(err).Str("feedId", feed.Id).Str("feedLink", feed.Link).Msg("error due to marshal feed")
				continue
			}

			id := asynq.TaskID(consts.TaskFeedImporter + feed.Id)

			if _, err = h.Client.EnqueueContext(ctx, asynq.NewTask(consts.TaskFeedImporter, payload), id); err != nil {
				if !errors.Is(err, asynq.ErrTaskIDConflict) {
					l.Error().Err(err).RawJSON("feed", payload).Msg("error due to enqueue feed")
				}
			}
		}

		if len(feeds) < int(size) {
			break
		}

		index += uint64(size)
	}
	return nil
}
