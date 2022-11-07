package storage

import (
	"context"
	"fmt"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ UserStorage = (*userStorage)(nil)

type UserStorage interface {
	Find(ctx context.Context, criteria mongoext.Criteria) ([]models.User, error)
	FindById(ctx context.Context, id string) (models.User, error)
	FindByUsername(ctx context.Context, username string) (models.User, error)
	Count(ctx context.Context, filter any) (int64, error)
	Save(ctx context.Context, model *models.User) error
	Delete(ctx context.Context, id string) error
}

type userStorage struct {
	c *mongo.Collection
}

func newUserStorage(db *mongo.Database) *userStorage {
	return &userStorage{c: db.Collection("users")}
}

func (s *userStorage) indexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, mongoext.Timeout)
	defer cancel()

	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"username", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"email", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"roles", 1}}},
		{Keys: bson.D{
			{"providers.name", 1},
			{"providers.value", 1},
			{"created_at", 1},
			{"updated_at", 1},
		}},
	}

	_, err := s.c.Indexes().CreateMany(ctx, indexes)

	return err
}

func (s *userStorage) Find(ctx context.Context, criteria mongoext.Criteria) ([]models.User, error) {
	return mongoext.FindByCriteria[models.User](ctx, s.c, criteria)
}

func (s *userStorage) FindById(ctx context.Context, id string) (models.User, error) {
	return mongoext.FindOne[models.User](ctx, s.c, bson.D{{"_id", id}})
}

func (s *userStorage) FindByUsername(ctx context.Context, username string) (user models.User, err error) {
	criteria := mongoext.Criteria{
		Size: 1,
		Filter: mongoext.Filter{mongoext.Or(
			mongoext.Eq("username", username),
			mongoext.Eq("email", username),
		)}.Build(),
	}
	items, err := s.Find(ctx, criteria)
	if err != nil {
		return
	}
	if len(items) != 1 {
		return user, fmt.Errorf(mongoext.ErrMsgQuery, mongo.ErrNoDocuments)
	}
	return items[0], nil
}

func (s *userStorage) Count(ctx context.Context, filter any) (int64, error) {
	return mongoext.Count(ctx, s.c, filter)
}

func (s *userStorage) Save(ctx context.Context, model *models.User) error {
	now := time.Now().UTC()
	model.UpdatedAt = now

	data, err := mongoext.ToBson(model)
	if err != nil {
		return err
	}

	delete(data, "_id")
	delete(data, "roles")
	delete(data, "providers")
	delete(data, "created_at")

	var write []mongo.WriteModel
	if deleteProviders := lo.Map(model.Providers, func(item models.ProviderData, _ int) models.Provider {
		return item.Name
	}); len(deleteProviders) > 0 {
		write = append(write, mongo.
			NewUpdateOneModel().
			SetFilter(bson.D{{"_id", model.Id}}).
			SetUpdate(bson.M{
				"$pull": bson.M{
					"providers": bson.M{
						"name": bson.M{"$in": deleteProviders},
					},
				},
			}),
		)
	}

	if len(model.DeleteRoles) > 0 {
		write = append(write, mongo.
			NewUpdateOneModel().
			SetFilter(bson.D{{"_id", model.Id}}).
			SetUpdate(bson.M{
				"$pull": bson.M{
					"roles": bson.M{"$in": model.DeleteRoles},
				},
			}),
		)
	}

	addProviders := lo.Filter(model.Providers, func(item models.ProviderData, _ int) bool {
		return !item.Delete
	})

	write = append(write, mongo.
		NewUpdateOneModel().
		SetUpsert(true).
		SetFilter(bson.D{{"_id", model.Id}}).
		SetUpdate(bson.M{
			"$set": data,
			"$setOnInsert": bson.M{
				"created_at": now,
			},
			"$addToSet": bson.M{
				"roles":     bson.M{"$each": model.Roles},
				"providers": bson.M{"$each": addProviders},
			},
		}),
	)

	result, err := mongoext.BulkWrite(ctx, s.c, write)
	if err != nil {
		return err
	}

	if result.UpsertedCount > 0 {
		model.CreatedAt = now
	}

	return nil
}

func (s *userStorage) Delete(ctx context.Context, id string) error {
	return mongoext.Delete(ctx, s.c, bson.D{{"_id", id}})
}
