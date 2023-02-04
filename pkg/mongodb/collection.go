package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Timeout = 5 * time.Second

const (
	ErrMsgDecode    = "failed to decode document due to error: %w"
	ErrMsgQuery     = "failed to execute query due to error: %w"
	ErrMsgMarshal   = "failed to marshal document due to error: %w"
	ErrMsgUnmarshal = "failed to unmarshal document due to error: %w"
)

func Save(ctx context.Context, c *mongo.Collection, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	opts = append(opts, options.Update().SetUpsert(true))

	if result, err := c.UpdateOne(ctx, filter, update, opts...); err != nil {
		return nil, fmt.Errorf(ErrMsgQuery, err)
	} else {
		return result, nil
	}
}

func SaveMany(ctx context.Context, c *mongo.Collection, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	opts = append(opts, options.Update().SetUpsert(true))

	if result, err := c.UpdateMany(ctx, filter, update, opts...); err != nil {
		return nil, fmt.Errorf(ErrMsgQuery, err)
	} else {
		return result, nil
	}
}

func BulkWrite(ctx context.Context, c *mongo.Collection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	if result, err := c.BulkWrite(ctx, models, opts...); err != nil {
		return nil, fmt.Errorf(ErrMsgQuery, err)
	} else {
		return result, nil
	}
}

func Remove(ctx context.Context, c *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	if result, err := c.DeleteOne(ctx, filter, opts...); err != nil {
		return fmt.Errorf(ErrMsgQuery, err)
	} else if result.DeletedCount == 0 {
		return fmt.Errorf(ErrMsgQuery, mongo.ErrNoDocuments)
	}
	return nil
}

func RemoveMany(ctx context.Context, c *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	if result, err := c.DeleteMany(ctx, filter, opts...); err != nil {
		return fmt.Errorf(ErrMsgQuery, err)
	} else if result.DeletedCount == 0 {
		return fmt.Errorf(ErrMsgQuery, mongo.ErrNoDocuments)
	}
	return nil
}

func FindOne[T any](ctx context.Context, c *mongo.Collection, filter any, opts ...*options.FindOneOptions) (doc T, err error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	err = DecodeOne(c.FindOne(ctx, filter, opts...), &doc)
	return
}

func Find[T any](ctx context.Context, c *mongo.Collection, filter any, opts ...*options.FindOptions) (docs []T, err error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	if result, err := c.Find(ctx, normalize(filter), opts...); err != nil {
		return nil, fmt.Errorf(ErrMsgQuery, err)
	} else {
		err = DecodeAll(ctx, result, &docs)
	}
	return
}

func Count(ctx context.Context, c *mongo.Collection, filter any, opts ...*options.CountOptions) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	count, err := c.CountDocuments(ctx, normalize(filter), opts...)
	if err != nil {
		return 0, fmt.Errorf(ErrMsgQuery, err)
	}
	return count, nil
}

func DecodeOne(r *mongo.SingleResult, doc any) error {
	if r.Err() != nil {
		return fmt.Errorf(ErrMsgQuery, r.Err())
	}
	if err := r.Decode(doc); err != nil {
		return fmt.Errorf(ErrMsgDecode, err)
	}
	return nil
}

func DecodeAll(ctx context.Context, cur *mongo.Cursor, docs any) error {
	if cur.Err() != nil {
		return fmt.Errorf(ErrMsgQuery, cur.Err())
	}
	if err := cur.All(ctx, docs); err != nil {
		return fmt.Errorf(ErrMsgDecode, err)
	}
	return nil
}

func ToBson(doc any) (bson.M, error) {
	data, err := bson.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf(ErrMsgMarshal, err)
	}

	var m bson.M
	if err = bson.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf(ErrMsgUnmarshal, err)
	}
	return m, nil
}

func Pagination(index int64, size int64) *options.FindOptions {
	o := options.Find()
	if index >= 0 {
		o.SetSkip(index)
	}
	if size > 0 {
		o.SetLimit(size)
	}
	return o
}

func normalize(filter any) any {
	if filter == nil {
		return bson.M{}
	}
	return filter
}
