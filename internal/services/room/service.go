package room

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"go.uber.org/zap"
)

type Service interface {
	ChatMemberUpdated(ctx context.Context, chat tgbotapi.Chat, deleted bool, permissions ...models.RoomPermission) error
}

func (p *Plugin) ChatMemberUpdated(ctx context.Context, chat tgbotapi.Chat, deleted bool, permissions ...models.RoomPermission) error {
	var err error
	var room models.Room
	if room, err = p.storage.FindById(ctx, chat.ID); err != nil {
		room.Id = chat.ID
	}

	room.Type = models.RoomType(chat.Type)
	room.Title = chat.Title
	room.UserName = chat.UserName
	room.FirstName = chat.FirstName
	room.LastName = chat.LastName
	room.SetDeleted(deleted).SetPermissions(permissions)

	if err = p.storage.Save(ctx, &room); err != nil {
		p.log.Warn("error due to save room", zap.Error(err), zap.Any("room", room))
		return err
	}

	p.sender.SendView(0, tgbotsender.ViewRoom, room)
	return nil
}
