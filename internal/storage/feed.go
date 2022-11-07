package storage

import (
	"context"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ FeedStorage = (*feedStorage)(nil)

type FeedStorage interface {
	Save(ctx context.Context, model *models.Feed) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	Find(ctx context.Context, criteria mongoext.Criteria) ([]models.Feed, error)
	FindByLink(ctx context.Context, link string) (models.Feed, error)
	FindById(ctx context.Context, id string) (models.Feed, error)
}

type feedStorage struct {
	c *mongo.Collection
}

func newFeedStorage(db *mongo.Database) *feedStorage {
	return &feedStorage{c: db.Collection("feeds")}
}

func (s *feedStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, mongoext.Timeout)
	defer cancel()

	data := []mongo.IndexModel{
		{Keys: bson.D{{"link", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"languages", 1},
			{"host", 1},
			{"enabled", 1},
			{"created_at", 1},
			{"updated_at", 1},
		}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, data)

	return err
}

func (s *feedStorage) Save(ctx context.Context, model *models.Feed) error {
	now := time.Now().UTC()
	model.UpdatedAt = now

	data, err := mongoext.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "by")
	delete(data, "created_at")

	result, err := mongoext.Save(ctx, s.c, bson.D{{"_id", model.Id}}, bson.M{
		"$set": data,
		"$setOnInsert": bson.M{
			"by":         model.By,
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

func (s *feedStorage) Delete(ctx context.Context, id string) error {
	return mongoext.Delete(ctx, s.c, bson.D{{"_id", id}})
}

func (s *feedStorage) Count(ctx context.Context, filter any) (int64, error) {
	return mongoext.Count(ctx, s.c, filter)
}

func (s *feedStorage) Find(ctx context.Context, criteria mongoext.Criteria) ([]models.Feed, error) {
	return mongoext.FindByCriteria[models.Feed](ctx, s.c, criteria)
}

func (s *feedStorage) FindByLink(ctx context.Context, link string) (models.Feed, error) {
	return mongoext.FindOne[models.Feed](ctx, s.c, bson.D{{"link", link}})
}

func (s *feedStorage) FindById(ctx context.Context, id string) (models.Feed, error) {
	return mongoext.FindOne[models.Feed](ctx, s.c, bson.D{{"_id", id}})
}
