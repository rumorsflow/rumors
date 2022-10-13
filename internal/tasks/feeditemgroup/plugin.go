package feeditemgroup

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	rumorscast "github.com/rumorsflow/rumors/internal/pkg/cast"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

const PluginName = consts.TaskFeedItemGroup

type Plugin struct {
	log         *zap.Logger
	client      *asynq.Client
	roomStorage storage.RoomStorage
}

func (p *Plugin) Init(log *zap.Logger, client *asynq.Client, roomStorage storage.RoomStorage) error {
	p.log = log
	p.client = client
	p.roomStorage = roomStorage
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var items []models.FeedItem
	if err := json.Unmarshal(task.Payload(), &items); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}

	if len(items) == 0 {
		p.log.Debug("items group are empty")
		return nil
	}

	var b strings.Builder
	b.WriteString(items[0].FeedId)
	group := map[string][]models.FeedItem{
		items[0].FeedId: {items[0]},
	}
	for _, item := range items[1:] {
		if _, ok := group[item.FeedId]; ok {
			group[item.FeedId] = append(group[item.FeedId], item)
		} else {
			group[item.FeedId] = []models.FeedItem{item}
			b.WriteString(",")
			b.WriteString(item.FeedId)
		}
	}

	size := 20
	query := make(url.Values)
	query.Set(mongoext.QuerySize, cast.ToString(size))
	query.Set("f[0][0][field]", "broadcast")
	query.Set("f[0][0][condition]", "in")
	query.Set("f[0][0][value]", b.String())
	query.Set("f[1][0][field]", "deleted")
	query.Set("f[1][0][value]", "false")

	for index := 0; ; index += size {
		query.Set(mongoext.QueryIndex, cast.ToString(index))

		rooms, err := p.roomStorage.Find(ctx, mongoext.C(query, "f"))
		if err != nil {
			p.log.Error("error due to find rooms", zap.Error(err))
			return nil
		}

		for _, room := range rooms {
			var roomItems []models.FeedItem
			for _, id := range *room.Broadcast {
				if data, ok := group[id]; ok {
					roomItems = append(roomItems, data...)
				}
			}
			if len(roomItems) > 0 {
				payload, _ := json.Marshal(roomItems)
				payload = append(payload, rumorscast.Int64ToBytes(room.Id)...)
				t := asynq.NewTask(consts.TaskRoomBroadcast, payload)
				q := asynq.Queue(consts.QueueFeedItems)

				if _, err = p.client.EnqueueContext(ctx, t, q); err != nil {
					p.log.Error(
						"error due to enqueue room broadcast",
						zap.Error(err),
						zap.String("task", t.Type()),
						zap.ByteString("payload", payload),
					)
				}
			}
		}

		if len(rooms) < size {
			break
		}
	}

	return nil
}
