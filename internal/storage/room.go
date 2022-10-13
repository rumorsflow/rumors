package storage

import (
	"context"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ RoomStorage = (*roomStorage)(nil)

type RoomStorage interface {
	Save(ctx context.Context, model *models.Room) error
	Count(ctx context.Context, filter any) (int64, error)
	Find(ctx context.Context, criteria mongoext.Criteria) ([]models.Room, error)
	FindById(ctx context.Context, id int64) (models.Room, error)
}

type roomStorage struct {
	c *mongo.Collection
}

func newRoomStorage(db *mongo.Database) *roomStorage {
	return &roomStorage{c: db.Collection("rooms")}
}

func (s *roomStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, mongoext.Timeout)
	defer cancel()

	_, err := s.c.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"broadcast", 1},
		{"deleted", 1},
		{"created_at", 1},
		{"updated_at", 1},
	}})

	return err
}

func (s *roomStorage) Save(ctx context.Context, model *models.Room) error {
	now := time.Now().UTC()
	model.UpdatedAt = now

	data, err := mongoext.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "created_at")

	result, err := mongoext.Save(ctx, s.c, bson.D{{"_id", model.Id}}, bson.M{
		"$set": data,
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	})
	if err != nil {
		return err
	}

	if result.UpsertedCount > 0 {
		model.CreatedAt = now
	}

	return nil
}

func (s *roomStorage) Count(ctx context.Context, filter any) (int64, error) {
	return mongoext.Count(ctx, s.c, filter)
}

func (s *roomStorage) Find(ctx context.Context, criteria mongoext.Criteria) ([]models.Room, error) {
	return mongoext.FindByCriteria[models.Room](ctx, s.c, criteria)
}

func (s *roomStorage) FindById(ctx context.Context, id int64) (models.Room, error) {
	return mongoext.FindOne[models.Room](ctx, s.c, bson.D{{"_id", id}})
}
