package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const ErrMsgDatabase = "failed to create mongodb database due to error: %w"

var ErrNoDB = errors.New("database name not found in URI")

func GetDB(ctx context.Context, uri string) (*mongo.Database, error) {
	dbName, err := GetDBName(uri)
	if err != nil {
		return nil, fmt.Errorf(ErrMsgDatabase, err)
	}

	client, err := GetClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}

func GetDBName(uri string) (string, error) {
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return "", err
	}
	if len(cs.Database) == 0 {
		return "", ErrNoDB
	}
	return cs.Database, nil
}
