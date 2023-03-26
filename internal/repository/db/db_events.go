package db

import (
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

func BeforeSave[T repository.Entity](entity T) (bson.M, error) {
	now := time.Now()
	created := now

	rv := reflect.ValueOf(entity).Elem()

	cField := rv.FieldByName("CreatedAt")
	if !cField.IsZero() {
		created = cField.Interface().(time.Time)
	}

	uField := rv.FieldByName("UpdatedAt")
	if uField.CanSet() {
		uField.Set(reflect.ValueOf(now))
	}

	data, err := mongodb.ToBson(entity)
	if err != nil {
		return nil, err
	}

	delete(data, "_id")
	delete(data, "created_at")

	return bson.M{
		"$set": data,
		"$setOnInsert": bson.M{
			"created_at": created,
		},
	}, nil
}

func AfterSave[T repository.Entity](entity T, result *mongo.UpdateResult) error {
	if result.UpsertedCount > 0 {
		rv := reflect.ValueOf(entity).Elem()
		fCreatedAt := rv.FieldByName("CreatedAt")
		fUpdatedAt := rv.FieldByName("UpdatedAt")

		if !fUpdatedAt.IsZero() && fCreatedAt.IsZero() && fCreatedAt.CanSet() {
			fCreatedAt.Set(fUpdatedAt)
		}
	}
	return nil
}

func ArticleBeforeSave(entity *entity.Article) (bson.M, error) {
	data, err := BeforeSave(entity)
	if err != nil {
		return data, err
	}

	set := data["$set"].(bson.M)
	delete(set, "site_id")
	delete(set, "source")
	delete(set, "link")
	delete(set, "pub_date")

	insert := data["$setOnInsert"].(bson.M)
	insert["site_id"] = entity.SiteID
	insert["source"] = entity.Source
	insert["link"] = entity.Link
	insert["pub_date"] = entity.PubDate

	return data, err
}

func ChatBeforeSave(entity *entity.Chat) (bson.M, error) {
	data, err := BeforeSave(entity)
	if err != nil {
		return data, err
	}

	data["$setOnInsert"].(bson.M)["telegram_id"] = entity.TelegramID
	delete(data["$set"].(bson.M), "telegram_id")

	return data, err
}
