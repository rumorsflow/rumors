package subscribed

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

const PluginName = consts.TaskSubscribed

type Plugin struct {
	log         *zap.Logger
	feedStorage storage.FeedStorage
	roomStorage storage.RoomStorage
	sender      tgbotsender.TelegramSender
}

func (p *Plugin) Init(
	log *zap.Logger,
	feedStorage storage.FeedStorage,
	roomStorage storage.RoomStorage,
	sender tgbotsender.TelegramSender,
) error {
	p.log = log
	p.feedStorage = feedStorage
	p.roomStorage = roomStorage
	p.sender = sender
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

	room, err := p.roomStorage.FindById(ctx, message.Chat.ID)
	if err != nil {
		p.log.Error("error due to find room", zap.Error(err))
		return err
	}

	var sources []string
	if room.Broadcast != nil && len(*room.Broadcast) > 0 {
		q := make(url.Values)
		q.Set(mongoext.QuerySize, "100")
		q.Set("f[0][0][field]", "_id")
		q.Set("f[0][0][condition]", mongoext.CondIn)
		q.Set("f[0][0][value]", strings.Join(*room.Broadcast, ","))
		q.Set("f[1][0][field]", "enabled")
		q.Set("f[1][0][value]", "true")

		var items []models.Feed
		for index := 0; ; index += 100 {
			q.Set(mongoext.QueryIndex, strconv.Itoa(index))
			criteria := mongoext.C(q, "f")

			items, err = p.feedStorage.Find(ctx, criteria)
			if err != nil {
				p.log.Error("error due to find feeds", zap.Error(err))
				p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, consts.ErrMsgTryLater)
				return nil
			}

			sources = append(sources, lo.Map(items, func(item models.Feed, _ int) string {
				return item.Host
			})...)

			if len(items) <= 100 {
				break
			}
		}

		sources = lo.Uniq(sources)
	}

	p.sender.SendView(message.Chat.ID, tgbotsender.ViewSubscribed, sources)

	return nil
}
