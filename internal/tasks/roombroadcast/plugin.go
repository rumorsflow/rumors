package roombroadcast

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/pkg/cast"
	"github.com/rumorsflow/rumors/internal/pkg/str"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"go.uber.org/zap"
	"strings"
)

const PluginName = consts.TaskRoomBroadcast

type Plugin struct {
	log    *zap.Logger
	sender tgbotsender.TelegramSender
}

func (p *Plugin) Init(log *zap.Logger, sender tgbotsender.TelegramSender) error {
	p.log = log
	p.sender = sender
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(_ context.Context, task *asynq.Task) error {
	if len(task.Payload()) <= 8 {
		return nil
	}

	var items []models.FeedItem
	if err := json.Unmarshal(task.Payload()[:len(task.Payload())-8], &items); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}

	chatId := cast.BytesToInt64(task.Payload()[len(task.Payload())-8:])
	group := make(map[string][]models.FeedItem)

	for _, item := range items {
		d := item.Domain()
		if item.Desc != nil {
			desc := str.MaxLen(*item.Desc, 500)
			desc = strings.TrimRight(desc, ".")
			desc = fmt.Sprintf("%s…", desc)
			item.Desc = &desc
		}
		if _, ok := group[d]; !ok {
			group[d] = []models.FeedItem{item}
		} else {
			group[d] = append(group[d], item)
		}
	}

	p.sender.SendView(chatId, tgbotsender.ViewFeedItems, group)

	return nil
}
