package room

import (
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"go.uber.org/zap"
)

const PluginName = "room_service"

type Plugin struct {
	log         *zap.Logger
	roomStorage storage.RoomStorage
	feedStorage storage.FeedStorage
	sender      tgbotsender.TelegramSender
}

func (p *Plugin) Init(
	log *zap.Logger,
	roomStorage storage.RoomStorage,
	feedStorage storage.FeedStorage,
	sender tgbotsender.TelegramSender,
) error {
	p.log = log
	p.roomStorage = roomStorage
	p.feedStorage = feedStorage
	p.sender = sender
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
