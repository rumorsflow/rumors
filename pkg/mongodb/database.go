package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var ErrNoDB = errors.New("database name not found in URI")

const (
	ErrMsgClient   = "failed to create mongodb client due to error: %w"
	ErrMsgDatabase = "failed to create mongodb database due to error: %w"
)

type Database struct {
	*mongo.Database
}

func NewDatabase(ctx context.Context, cfg *Config) (*Database, error) {
	dbName, err := ExtractDatabaseName(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf(ErrMsgDatabase, err)
	}

	client, err := NewClient(ctx, cfg.URI)
	if err != nil {
		return nil, err
	}

	db := &Database{Database: client.Database(dbName)}

	if cfg.Ping {
		if err = db.Ping(ctx); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *Database) Ping(ctx context.Context) error {
	if err := db.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not connect to MongoDB: %w", err)
	}
	return nil
}

func (db *Database) Close(ctx context.Context) error {
	return db.Client().Disconnect(ctx)
}

func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf(ErrMsgClient, err)
	}
	return client, nil
}

func ExtractDatabaseName(uri string) (string, error) {
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return "", err
	}
	if len(cs.Database) == 0 {
		return "", ErrNoDB
	}
	return cs.Database, nil
}
