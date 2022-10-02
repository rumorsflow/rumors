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

type FilterFeedItems struct {
	FeedIds  []string `json:"feed_ids,omitempty" query:"feed_ids[]"`
	Link     *string  `json:"link,omitempty" query:"link"`
	Author   *string  `json:"author,omitempty" query:"author"`
	Category *string  `json:"category,omitempty" query:"category"`
}

func (f *FilterFeedItems) SetIds(feedIds ...string) *FilterFeedItems {
	f.FeedIds = feedIds
	return f
}

func (f *FilterFeedItems) SetLink(link string) *FilterFeedItems {
	f.Link = &link
	return f
}

func (f *FilterFeedItems) SetAuthor(author string) *FilterFeedItems {
	f.Author = &author
	return f
}

func (f *FilterFeedItems) SetCategory(category string) *FilterFeedItems {
	f.Category = &category
	return f
}

func (f *FilterFeedItems) build() any {
	if f == nil {
		return bson.D{}
	}

	var filter mongodb.Filter

	if len(f.FeedIds) > 0 {
		filter = append(filter, mongodb.In("feed_id", slice.ToAny(f.FeedIds)...))
	}

	if f.Link != nil {
		filter = append(filter, mongodb.Regex("link", *f.Link, "i"))
	}

	if f.Author != nil {
		filter = append(filter, mongodb.Regex("authors", *f.Author, "i"))
	}

	if f.Category != nil {
		filter = append(filter, mongodb.Regex("categories", *f.Category, "i"))
	}

	return filter.Build()
}

type FeedItemStorage interface {
	Find(ctx context.Context, filter *FilterFeedItems, index uint64, size uint32) ([]models.FeedItem, error)
	Save(ctx context.Context, model models.FeedItem) error
}

type feedItemStorage struct {
	c *mongo.Collection
}

func NewFeedItemStorage(ctx context.Context, db *mongo.Database) (FeedItemStorage, error) {
	s := &feedItemStorage{
		c: db.Collection("feed_items"),
	}
	if err := s.indexes(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *feedItemStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data := []mongo.IndexModel{
		{Keys: bson.D{
			{"guid", 1},
			{"pub_date", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"feed_id", 1}, {"link", 1}, {"pub_date", 1}, {"created_at", 1}}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, data)

	return err
}

func (s *feedItemStorage) Find(ctx context.Context, filter *FilterFeedItems, index uint64, size uint32) ([]models.FeedItem, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := mongodb.FindOptions(index, size).SetSort(bson.D{{"pub_date", -1}})
	cur, err := s.c.Find(ctx, filter.build(), opts)
	if err != nil {
		return nil, fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	return mongodb.DecodeAll[models.FeedItem](ctx, cur)
}

func (s *feedItemStorage) Save(ctx context.Context, model models.FeedItem) error {
	data, err := mongodb.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "feed_id")
	delete(data, "guid")
	delete(data, "pub_date")
	delete(data, "created_at")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err = s.c.UpdateOne(
		ctx,
		bson.M{"_id": model.Id},
		bson.M{
			"$set": data,
			"$setOnInsert": bson.M{
				"feed_id":    model.FeedId,
				"guid":       model.Guid,
				"pub_date":   model.PubDate,
				"created_at": model.CreatedAt,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return fmt.Errorf(fmt.Sprintf("feed save: %s", mongodb.ErrMsgQuery), err)
	}

	return nil
}
