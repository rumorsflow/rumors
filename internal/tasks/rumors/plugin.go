package rumors

import (
	"context"
	"encoding/json"
	"github.com/go-fc/slice"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	rumorscast "github.com/rumorsflow/rumors/internal/pkg/cast"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"unicode/utf8"
)

const PluginName = consts.TaskRumors

type Plugin struct {
	log             *zap.Logger
	client          *asynq.Client
	feedStorage     storage.FeedStorage
	feedItemStorage storage.FeedItemStorage
}

var feedsCriteria mongoext.Criteria

func init() {
	q := make(url.Values)
	q.Set(mongoext.QueryIndex, "0")
	q.Set(mongoext.QuerySize, "100")
	q.Set("f[0][0][field]", "enabled")
	q.Set("f[0][0][value]", "true")

	feedsCriteria = mongoext.C(q, "f")
}

func (p *Plugin) Init(
	log *zap.Logger,
	client *asynq.Client,
	feedStorage storage.FeedStorage,
	feedItemStorage storage.FeedItemStorage,
) error {
	p.log = log
	p.client = client
	p.feedStorage = feedStorage
	p.feedItemStorage = feedItemStorage
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var message tgbotapi.Message
	if err := json.Unmarshal(task.Payload(), &message); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}

	feeds, err := p.feedStorage.Find(ctx, feedsCriteria)
	if err != nil {
		p.log.Error("error due to find feeds", zap.Error(err))
		return err
	}

	feedIds := lo.Map(feeds, func(item models.Feed, _ int) string {
		return item.Id
	})

	rf := rumorsFilter(message.CommandArguments(), feedIds)
	items, err := p.feedItemStorage.Find(ctx, mongoext.C(rf, "f"))
	if err != nil {
		p.log.Error("error due to find feed items", zap.Error(err))
		return err
	}

	payload, _ := json.Marshal(items)
	payload = append(payload, rumorscast.Int64ToBytes(message.Chat.ID)...)
	t := asynq.NewTask(consts.TaskRoomBroadcast, payload)

	if _, err = p.client.EnqueueContext(ctx, t); err != nil {
		p.log.Error(
			"error due to enqueue room broadcast",
			zap.Error(err),
			zap.String("task", t.Type()),
			zap.ByteString("payload", payload),
		)
		return err
	}
	return nil
}

func rumorsFilter(args string, feedIds []string) url.Values {
	index, size, search := pagination(args)
	q := make(url.Values)
	q.Set(mongoext.QueryIndex, cast.ToString(index))
	q.Set(mongoext.QuerySize, cast.ToString(size))
	q.Set(mongoext.QuerySortArr, "-pub_date")
	q.Set("f[0][0][field]", "feed_id")
	q.Set("f[0][0][condition]", "in")
	q.Set("f[0][0][value]", strings.Join(feedIds, ","))

	if utf8.RuneCountInString(search) > 0 {
		q.Set("f[1][0][field]", "link")
		q.Set("f[1][0][condition]", "regex")
		q.Set("f[1][0][value]", search)
		q.Set("f[1][1][field]", "title")
		q.Set("f[1][1][condition]", "regex")
		q.Set("f[1][1][value]", search)
		q.Set("f[1][2][field]", "desc")
		q.Set("f[1][2][condition]", "regex")
		q.Set("f[1][2][value]", search)
		q.Set("f[1][3][field]", "categories")
		q.Set("f[1][3][condition]", "regex")
		q.Set("f[1][3][value]", search)
	}
	return q
}

func pagination(args string) (i uint64, s uint64, search string) {
	a := strings.Fields(args)
	i = cast.ToUint64(slice.Safe(a, 0))
	if s = cast.ToUint64(slice.Safe(a, 1)); s > 0 {
		if s > 20 {
			s = 20
		}
	} else {
		s = 10
	}

	if len(a) > 2 {
		search = strings.TrimSpace(strings.Join(a[2:], " "))
	}
	return
}
