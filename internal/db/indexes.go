package db

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SiteIndexes(indexView mongo.IndexView) error {
	if _, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{"domain", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"languages", 1},
			{"title", 1},
			{"enabled", 1},
		}},
	}); err != nil {
		return fmt.Errorf("%s %w", repository.OpIndexes, err)
	}
	return nil
}

func ArticleIndexes(indexView mongo.IndexView) error {
	if _, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{"link", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"pub_date", 1}}},
		{Keys: bson.D{{"created_at", 1}}},
		{Keys: bson.D{{"updated_at", 1}}},
		{Keys: bson.D{{"site_id", 1}, {"lang", 1}}},
		{Keys: bson.D{
			{"source", 1},
			{"lang", 1},
			{"categories", 1},
		}},
	}); err != nil {
		return fmt.Errorf("%s %w", repository.OpIndexes, err)
	}
	return nil
}

func ChatIndexes(indexView mongo.IndexView) error {
	if _, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{"telegram_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"broadcast", 1},
			{"blocked", 1},
			{"deleted", 1},
			{"created_at", 1},
			{"updated_at", 1},
		}},
	}); err != nil {
		return fmt.Errorf("%s %w", repository.OpIndexes, err)
	}
	return nil
}

func JobIndexes(indexView mongo.IndexView) error {
	if _, err := indexView.CreateOne(context.Background(), mongo.IndexModel{Keys: bson.D{
		{"name", 1},
		{"enabled", 1},
		{"created_at", 1},
		{"updated_at", 1},
	}}); err != nil {
		return fmt.Errorf("%s %w", repository.OpIndexes, err)
	}
	return nil
}

func SysUserIndexes(indexView mongo.IndexView) error {
	if _, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{"username", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"email", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"created_at", 1},
			{"updated_at", 1},
		}},
	}); err != nil {
		return fmt.Errorf("%s %w", repository.OpIndexes, err)
	}
	return nil
}
