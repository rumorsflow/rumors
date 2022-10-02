package storage

import (
	"context"
	"fmt"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/mongodb"
	"github.com/iagapie/rumors/pkg/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FilterRooms struct {
	Ids       []string `json:"ids,omitempty" query:"ids[]"`
	ChatIds   []int64  `json:"chat_ids,omitempty" query:"chat_ids[]"`
	Title     *string  `json:"title,omitempty" query:"title"`
	Broadcast *bool    `json:"broadcast,omitempty" query:"broadcast"`
	Deleted   *bool    `json:"deleted,omitempty" query:"deleted"`
}

func (f *FilterRooms) SetIds(ids ...string) *FilterRooms {
	f.Ids = ids
	return f
}

func (f *FilterRooms) SetChatIds(chatIds ...int64) *FilterRooms {
	f.ChatIds = chatIds
	return f
}

func (f *FilterRooms) SetTitle(title string) *FilterRooms {
	f.Title = &title
	return f
}

func (f *FilterRooms) SetBroadcast(broadcast bool) *FilterRooms {
	f.Broadcast = &broadcast
	return f
}

func (f *FilterRooms) SetDeleted(deleted bool) *FilterRooms {
	f.Deleted = &deleted
	return f
}

func (f *FilterRooms) build() any {
	if f == nil {
		return bson.D{}
	}

	var filter mongodb.Filter

	if len(f.Ids) > 0 {
		filter = append(filter, mongodb.In("_id", slice.ToAny(f.Ids)...))
	}

	if len(f.ChatIds) > 0 {
		filter = append(filter, mongodb.In("chat_id", slice.ToAny(f.ChatIds)...))
	}

	if f.Title != nil {
		filter = append(filter, mongodb.Eq("title", *f.Title))
	}

	if f.Broadcast != nil {
		filter = append(filter, mongodb.Eq("broadcast", *f.Broadcast))
	}

	if f.Deleted != nil {
		filter = append(filter, mongodb.Eq("deleted", *f.Deleted))
	}

	return filter.Build()
}

type RoomStorage interface {
	Save(ctx context.Context, model models.Room) error
	Find(ctx context.Context, filter *FilterRooms, index uint64, size uint32) ([]models.Room, error)
	FindByChatId(ctx context.Context, chatId int64) (models.Room, error)
	FindById(ctx context.Context, id string) (models.Room, error)
}

type roomStorage struct {
	c *mongo.Collection
}

func NewRoomStorage(ctx context.Context, db *mongo.Database) (RoomStorage, error) {
	s := &roomStorage{
		c: db.Collection("rooms"),
	}
	if err := s.indexes(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *roomStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data := []mongo.IndexModel{
		{Keys: bson.D{{"chat_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"broadcast", 1}, {"deleted", 1}, {"created_at", 1}}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, data)

	return err
}

func (s *roomStorage) Save(ctx context.Context, model models.Room) error {
	data, err := mongodb.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "chat_id")
	delete(data, "created_at")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err = s.c.UpdateOne(
		ctx,
		bson.M{"_id": model.Id},
		bson.M{
			"$set": data,
			"$setOnInsert": bson.M{
				"chat_id":    model.ChatId,
				"created_at": model.CreatedAt,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return fmt.Errorf(fmt.Sprintf("room save: %s", mongodb.ErrMsgQuery), err)
	}

	return nil
}

func (s *roomStorage) Find(ctx context.Context, filter *FilterRooms, index uint64, size uint32) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := mongodb.FindOptions(index, size).SetSort(bson.D{{"created_at", 1}})
	cur, err := s.c.Find(ctx, filter.build(), opts)
	if err != nil {
		return nil, fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	return mongodb.DecodeAll[models.Room](ctx, cur)
}

func (s *roomStorage) FindByChatId(ctx context.Context, chatId int64) (models.Room, error) {
	return s.findOne(ctx, new(FilterRooms).SetChatIds(chatId))
}

func (s *roomStorage) FindById(ctx context.Context, id string) (models.Room, error) {
	return s.findOne(ctx, new(FilterRooms).SetIds(id))
}

func (s *roomStorage) findOne(ctx context.Context, filter *FilterRooms) (models.Room, error) {
	data, err := s.Find(ctx, filter, 0, 1)
	if err != nil {
		return models.Room{}, err
	}
	if len(data) == 0 {
		return models.Room{}, fmt.Errorf("room error: %w", mongo.ErrNoDocuments)
	}
	return data[0], nil
}
