package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_ repository.ReadRepository[repository.Entity]      = (*Repository[repository.Entity])(nil)
	_ repository.WriteRepository[repository.Entity]     = (*Repository[repository.Entity])(nil)
	_ repository.ReadWriteRepository[repository.Entity] = (*Repository[repository.Entity])(nil)
)

var (
	ErrMissingDB         = errs.New("missing *mongo.Database")
	ErrMissingCollection = errs.New("missing collection name")
)

type Option[T repository.Entity] func(*Repository[T]) error

type Repository[T repository.Entity] struct {
	collection    *mongo.Collection
	entityFactory repository.EntityFactory[T]
	afterFind     func(entity T) error
	beforeSave    func(entity T) (bson.M, error)
	afterSave     func(entity T, result *mongo.UpdateResult) error
}

func WithEntityFactory[T repository.Entity](entityFactory repository.EntityFactory[T]) Option[T] {
	return func(r *Repository[T]) error {
		r.entityFactory = entityFactory
		return nil
	}
}

func WithAfterFind[T repository.Entity](afterFind func(entity T) error) Option[T] {
	return func(r *Repository[T]) error {
		r.afterFind = afterFind
		return nil
	}
}

func WithBeforeSave[T repository.Entity](beforeSave func(entity T) (bson.M, error)) Option[T] {
	return func(r *Repository[T]) error {
		r.beforeSave = beforeSave
		return nil
	}
}

func WithAfterSave[T repository.Entity](afterSave func(entity T, result *mongo.UpdateResult) error) Option[T] {
	return func(r *Repository[T]) error {
		r.afterSave = afterSave
		return nil
	}
}

func WithIndexes[T repository.Entity](indexes func(indexView mongo.IndexView) error) Option[T] {
	return func(r *Repository[T]) error {
		return indexes(r.collection.Indexes())
	}
}

func NewRepository[T repository.Entity](database *mongodb.Database, collection string, options ...Option[T]) (*Repository[T], error) {
	if database == nil {
		return nil, errs.E(repository.OpNew, ErrMissingDB)
	}
	if collection == "" {
		return nil, errs.E(repository.OpNew, ErrMissingCollection)
	}

	r := &Repository[T]{collection: database.Collection(collection)}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, errs.Errorf(repository.OpNew, "error while applying option: %w", err)
		}
	}

	return r, nil
}

func (r *Repository[T]) Count(ctx context.Context, filter any) (count int64, err error) {
	count, err = mongodb.Count(ctx, r.collection, filter)
	if err != nil {
		err = errs.E(repository.OpCount, err)
	}
	return
}

func (r *Repository[T]) Find(ctx context.Context, criteria *repository.Criteria) ([]T, error) {
	ctx, cancel := context.WithTimeout(ctx, mongodb.Timeout)
	defer cancel()

	iter, err := r.FindIter(ctx, criteria)
	if err != nil {
		return nil, errs.E(repository.OpFind, err)
	}

	var result []T

	for iter.Next(ctx) {
		result = append(result, iter.Entity())
	}

	if err = iter.Close(ctx); err != nil {
		return nil, errs.E(repository.OpFind, err)
	}

	return result, nil
}

func (r *Repository[T]) FindIter(ctx context.Context, criteria *repository.Criteria) (repository.Iter[T], error) {
	if r.entityFactory == nil {
		return nil, errs.E(repository.OpFindIter, repository.ErrMissingEntityFactory)
	}

	var filter any
	o := options.Find()

	if criteria != nil {
		filter = criteria.Filter
		o.Skip = criteria.Index
		o.Limit = criteria.Size
		o.Sort = criteria.Sort
	}

	cursor, err := r.collection.Find(ctx, filter, o)
	if err != nil {
		return nil, errs.Errorf(repository.OpFindIter, mongodb.ErrMsgQuery, err)
	}

	return &repository.Iterator[T]{
		Cursor:    cursor,
		Factory:   r.entityFactory,
		AfterFind: r.afterFind,
	}, nil
}

func (r *Repository[T]) FindByID(ctx context.Context, id uuid.UUID) (value T, err error) {
	if r.entityFactory == nil {
		return value, errs.E(repository.OpFindByID, repository.ErrMissingEntityFactory)
	}

	ctx, cancel := context.WithTimeout(ctx, mongodb.Timeout)
	defer cancel()

	result := r.collection.FindOne(ctx, bson.M{"_id": id.String()})
	entity := r.entityFactory.NewEntity()

	if err = mongodb.DecodeOne(result, entity); err != nil {
		return value, toRepoErr(repository.OpFindByID, err, id)
	}

	if r.afterFind != nil {
		if err = r.afterFind(entity); err != nil {
			return value, errs.E(repository.OpFindByID, id, fmt.Errorf(repository.ErrMsgAfterFind, err))
		}
	}

	return entity, nil
}

func (r *Repository[T]) Save(ctx context.Context, entity T) (err error) {
	id := entity.EntityID()
	if id == uuid.Nil {
		return errs.E(repository.OpSave, repository.ErrMissingEntityID)
	}

	var update bson.M

	if r.beforeSave == nil {
		update = bson.M{"$set": entity}
	} else if update, err = r.beforeSave(entity); err != nil {
		return toRepoErr(repository.OpSave, fmt.Errorf("failed before save due to error: %w", err), id)
	}

	var result *mongo.UpdateResult

	if result, err = mongodb.Save(ctx, r.collection, bson.M{"_id": id.String()}, update); err != nil {
		return toRepoErr(repository.OpSave, err, id)
	} else if r.afterSave != nil {
		if err = r.afterSave(entity, result); err != nil {
			return toRepoErr(repository.OpSave, fmt.Errorf("failed after save due to error: %w", err), id)
		}
	}

	return nil
}

func (r *Repository[T]) Remove(ctx context.Context, id uuid.UUID) error {
	if err := mongodb.Remove(ctx, r.collection, bson.M{"_id": id.String()}); err != nil {
		return toRepoErr(repository.OpRemove, err, id)
	}
	return nil
}

func toRepoErr(op errs.Op, err error, id uuid.UUID) error {
	if errs.Is(err, mongo.ErrNoDocuments) {
		return errs.E(op, id, fmt.Errorf(mongodb.ErrMsgQuery, repository.ErrEntityNotFound))
	}
	if mongo.IsDuplicateKeyError(err) {
		return errs.E(op, id, fmt.Errorf(mongodb.ErrMsgQuery, repository.ErrDuplicateKey))
	}
	return errs.E(op, id, err)
}
