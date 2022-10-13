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

var _ FeedItemStorage = (*feedItemStorage)(nil)

type FeedItemStorage interface {
	Save(ctx context.Context, model *models.FeedItem) error
	Count(ctx context.Context, filter any) (int64, error)
	Find(ctx context.Context, criteria mongoext.Criteria) ([]models.FeedItem, error)
	FindById(ctx context.Context, id string) (models.FeedItem, error)
}

type feedItemStorage struct {
	c *mongo.Collection
}

func newFeedItemStorage(db *mongo.Database) *feedItemStorage {
	return &feedItemStorage{c: db.Collection("feed_items")}
}

func (s *feedItemStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, mongoext.Timeout)
	defer cancel()

	data := []mongo.IndexModel{
		{Keys: bson.D{{"guid", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"feed_id", 1},
			{"link", 1},
			{"pub_date", 1},
			{"created_at", 1},
			{"updated_at", 1},
		}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, data)

	return err
}

func (s *feedItemStorage) Save(ctx context.Context, model *models.FeedItem) error {
	now := time.Now().UTC()
	model.UpdatedAt = now

	if model.PubDate.IsZero() {
		model.PubDate = now
	}

	data, err := mongoext.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "feed_id")
	delete(data, "guid")
	delete(data, "pub_date")
	delete(data, "created_at")

	result, err := mongoext.Save(ctx, s.c, bson.D{{"_id", model.Id}}, bson.M{
		"$set": data,
		"$setOnInsert": bson.M{
			"feed_id":    model.FeedId,
			"guid":       model.Guid,
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

func (s *feedItemStorage) Count(ctx context.Context, filter any) (int64, error) {
	return mongoext.Count(ctx, s.c, filter)
}

func (s *feedItemStorage) Find(ctx context.Context, criteria mongoext.Criteria) ([]models.FeedItem, error) {
	return mongoext.FindByCriteria[models.FeedItem](ctx, s.c, criteria)
}

func (s *feedItemStorage) FindById(ctx context.Context, id string) (models.FeedItem, error) {
	return mongoext.FindOne[models.FeedItem](ctx, s.c, bson.D{{"_id", id}})
}
