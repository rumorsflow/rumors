package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ErrMsgDecode    = "failed to decode document due to error: %w"
	ErrMsgQuery     = "failed to execute query due to error: %w"
	ErrMsgMarshal   = "failed to marshal document due to error: %w"
	ErrMsgUnmarshal = "failed to unmarshal document due to error: %w"
)

func FindOptions(index uint64, size uint32) *options.FindOptions {
	return options.Find().SetSkip(int64(index)).SetLimit(int64(size))
}

func DecodeOne[T any](r *mongo.SingleResult) (doc T, err error) {
	if r.Err() != nil {
		return doc, fmt.Errorf(ErrMsgQuery, r.Err())
	}
	if err = r.Decode(&doc); err != nil {
		return doc, fmt.Errorf(ErrMsgDecode, err)
	}
	return doc, nil
}

func DecodeAll[T any](ctx context.Context, cur *mongo.Cursor) (docs []T, err error) {
	if cur.Err() != nil {
		return docs, fmt.Errorf(ErrMsgQuery, cur.Err())
	}
	if err = cur.All(ctx, &docs); err != nil {
		return docs, fmt.Errorf(ErrMsgDecode, err)
	}
	return docs, nil
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
