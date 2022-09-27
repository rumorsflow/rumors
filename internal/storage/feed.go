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
	"time"
)

const timeout = 5 * time.Second

type FilterFeeds struct {
	Ids     []string `json:"ids,omitempty" query:"ids[]"`
	By      *int64   `json:"by,omitempty" query:"by"`
	Host    *string  `json:"host,omitempty" query:"host" validate:"omitempty,fqdn"`
	Link    *string  `json:"link,omitempty" query:"link" validate:"omitempty,url"`
	Enabled *bool    `json:"enabled,omitempty" query:"enabled"`
}

func (f *FilterFeeds) build() any {
	var filter mongodb.Filter

	if len(f.Ids) > 0 {
		filter = append(filter, mongodb.In("_id", slice.ToAny(f.Ids)...))
	}

	if f.By != nil {
		filter = append(filter, mongodb.Eq("by", *f.By))
	}

	if f.Host != nil {
		filter = append(filter, mongodb.Regex("host", *f.Host, "i"))
	}

	if f.Link != nil {
		filter = append(filter, mongodb.Eq("link", *f.Link))
	}

	if f.Enabled != nil {
		filter = append(filter, mongodb.Eq("enabled", *f.Enabled))
	}

	return filter.Build()
}

type FeedStorage interface {
	FindByLink(ctx context.Context, link string) (models.Feed, error)
	FindById(ctx context.Context, id string) (models.Feed, error)
	Find(ctx context.Context, filter FilterFeeds, index uint64, size uint32) ([]models.Feed, error)
	Save(ctx context.Context, model models.Feed) error
}

type feedStorage struct {
	c *mongo.Collection
}

func NewFeedStorage(ctx context.Context, db *mongo.Database) (FeedStorage, error) {
	s := &feedStorage{
		c: db.Collection("feeds"),
	}
	if err := s.indexes(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *feedStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data := []mongo.IndexModel{
		{Keys: bson.D{{"link", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"host", 1}, {"enabled", 1}, {"created_at", 1}}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, data)

	return err
}

func (s *feedStorage) Find(ctx context.Context, filter FilterFeeds, index uint64, size uint32) ([]models.Feed, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := mongodb.FindOptions(index, size).SetSort(bson.D{{"created_at", 1}})
	cur, err := s.c.Find(ctx, filter.build(), opts)
	if err != nil {
		return nil, fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	return mongodb.DecodeAll[models.Feed](ctx, cur)
}

func (s *feedStorage) FindByLink(ctx context.Context, link string) (models.Feed, error) {
	data, err := s.Find(ctx, FilterFeeds{Link: &link}, 0, 1)
	if err != nil {
		return models.Feed{}, err
	}
	if len(data) == 0 {
		return models.Feed{}, fmt.Errorf("feed error: %w", mongo.ErrNoDocuments)
	}
	return data[0], nil
}

func (s *feedStorage) FindById(ctx context.Context, id string) (models.Feed, error) {
	data, err := s.Find(ctx, FilterFeeds{Ids: []string{id}}, 0, 1)
	if err != nil {
		return models.Feed{}, err
	}
	if len(data) == 0 {
		return models.Feed{}, fmt.Errorf("feed error: %w", mongo.ErrNoDocuments)
	}
	return data[0], nil
}

func (s *feedStorage) Save(ctx context.Context, model models.Feed) error {
	data, err := mongodb.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "by")
	delete(data, "host")
	delete(data, "link")
	delete(data, "created_at")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err = s.c.UpdateOne(
		ctx,
		bson.M{"_id": model.Id},
		bson.M{
			"$set": data,
			"$setOnInsert": bson.M{
				"by":         model.By,
				"host":       model.Host,
				"link":       model.Link,
				"created_at": model.CreatedAt,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return fmt.Errorf(fmt.Sprintf("feed save: %s", mongodb.ErrMsgQuery), err)
	}

	return nil
}
