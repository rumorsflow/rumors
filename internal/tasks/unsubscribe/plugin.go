package unsubscribe

import (
	"context"
	"encoding/json"
	"fmt"
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
)

const PluginName = consts.TaskUnsubscribe

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

	if message.CommandArguments() == "" {
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, consts.ErrMsgRequiredSource)
		return nil
	}

	q := make(url.Values)
	q.Set(mongoext.QueryIndex, "0")
	q.Set(mongoext.QuerySize, "1000")
	q.Set("f[0][0][field]", "host")
	q.Set("f[0][0][value]", message.CommandArguments())
	criteria := mongoext.C(q, "f")

	items, err := p.feedStorage.Find(ctx, criteria)
	if err != nil {
		p.log.Error("error due to find feeds", zap.Error(err))
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, consts.ErrMsgTryLater)
		return nil
	}
	if len(items) == 0 {
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewNotFound, fmt.Sprintf(consts.ErrMsgNotFoundSource, message.CommandArguments()))
		return nil
	}

	room, err := p.roomStorage.FindById(ctx, message.Chat.ID)
	if err != nil {
		p.log.Error("error due to find room", zap.Error(err))
		return err
	}
	if room.Broadcast == nil || len(*room.Broadcast) == 0 {
		return nil
	}

	broadcast := lo.Without(*room.Broadcast, lo.Map(items, func(item models.Feed, _ int) string {
		return item.Id
	})...)
	room.Broadcast = &broadcast

	if err = p.roomStorage.Save(ctx, &room); err != nil {
		p.log.Error("error due to save room", zap.Error(err))
		p.sender.SendView(message.Chat.ID, tgbotsender.ViewError, consts.ErrMsgTryLater)
		return nil
	}

	p.sender.SendView(message.Chat.ID, tgbotsender.ViewSuccess, fmt.Sprintf(consts.SuccessMsgUnsubscribed, message.CommandArguments()))

	return nil
}
