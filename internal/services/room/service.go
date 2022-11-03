package room

import (
	"context"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"net/url"
)

var (
	ErrEmptyBroadcast = errors.New("room don't have broadcast feeds")
	ErrNotFoundFeeds  = errors.New("feeds not found")
)

type Service interface {
	ChatMemberUpdated(ctx context.Context, chat tgbotapi.Chat, deleted bool, permissions ...models.RoomPermission) error
	AddToBroadcastByHost(ctx context.Context, id int64, host string) error
	DeleteFromBroadcastByHost(ctx context.Context, id int64, host string) error
	AddToBroadcast(ctx context.Context, id int64, feeds []string) error
	DeleteFromBroadcast(ctx context.Context, id int64, feeds []string) error
}

func (p *Plugin) ChatMemberUpdated(ctx context.Context, chat tgbotapi.Chat, deleted bool, permissions ...models.RoomPermission) error {
	var err error
	var room models.Room
	if room, err = p.roomStorage.FindById(ctx, chat.ID); err != nil {
		room.Id = chat.ID
	}

	room.Type = models.RoomType(chat.Type)
	room.Title = chat.Title
	room.UserName = chat.UserName
	room.FirstName = chat.FirstName
	room.LastName = chat.LastName
	room.SetDeleted(deleted).SetPermissions(permissions)

	if err = p.roomStorage.Save(ctx, &room); err != nil {
		p.log.Warn("error due to save room", zap.Error(err), zap.Any("room", room))
		return err
	}

	p.sender.SendView(0, tgbotsender.ViewRoom, room)
	return nil
}

func (p *Plugin) AddToBroadcastByHost(ctx context.Context, id int64, host string) error {
	feeds, err := p.getFeedsByHost(ctx, host)
	if err != nil {
		return err
	}
	return p.AddToBroadcast(ctx, id, feeds)
}
func (p *Plugin) DeleteFromBroadcastByHost(ctx context.Context, id int64, host string) error {
	feeds, err := p.getFeedsByHost(ctx, host)
	if err != nil {
		return err
	}
	return p.DeleteFromBroadcast(ctx, id, feeds)
}

func (p *Plugin) AddToBroadcast(ctx context.Context, id int64, feeds []string) error {
	return p.change(ctx, id, func(room *models.Room) error {
		if room.Broadcast != nil {
			feeds = append(feeds, *room.Broadcast...)
		}
		room.SetBroadcast(lo.Uniq(feeds))
		return nil
	})
}

func (p *Plugin) DeleteFromBroadcast(ctx context.Context, id int64, feeds []string) error {
	return p.change(ctx, id, func(room *models.Room) error {
		if room.Broadcast == nil || len(*room.Broadcast) == 0 {
			return ErrEmptyBroadcast
		}
		room.SetBroadcast(lo.Without(*room.Broadcast, feeds...))
		return nil
	})
}

func (p *Plugin) change(ctx context.Context, id int64, change func(*models.Room) error) error {
	room, err := p.roomStorage.FindById(ctx, id)
	if err != nil {
		p.log.Warn("error due to find room", zap.Int64("room", id), zap.Error(err))
		return err
	}

	if err = change(&room); err != nil {
		p.log.Error("error due to change room", zap.Int64("room", id), zap.Error(err))
		return err
	}

	if err = p.roomStorage.Save(ctx, &room); err != nil {
		p.log.Error("error due to save room", zap.Int64("room", id), zap.Error(err))
		return err
	}

	return nil
}

func (p *Plugin) getFeedsByHost(ctx context.Context, host string) ([]string, error) {
	q := make(url.Values)
	q.Set(mongoext.QueryIndex, "0")
	q.Set(mongoext.QuerySize, "1000")
	q.Set("f[0][0][field]", "host")
	q.Set("f[0][0][value]", host)
	q.Set("f[1][0][field]", "enabled")
	q.Set("f[1][0][value]", "true")
	criteria := mongoext.C(q, "f")

	items, err := p.feedStorage.Find(ctx, criteria)
	if err != nil {
		p.log.Warn("error due to find feeds", zap.Error(err))
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrNotFoundFeeds
	}

	return lo.Map(items, func(item models.Feed, _ int) string {
		return item.Id
	}), nil
}
