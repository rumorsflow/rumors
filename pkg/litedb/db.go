package litedb

import (
	"fmt"
	"github.com/pocketbase/dbx"
	"os"
	"strings"
)

type DB struct {
	*dbx.DB
}

func New(path, name string) (*DB, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	file := fmt.Sprintf("%s/%s", strings.TrimRight(path, "/"), name)
	db, err := Connect(file)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (x *DB) Attach(path, file, name string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	f := fmt.Sprintf("%s/%s", strings.TrimRight(path, "/"), file)
	_, err := x.DB.NewQuery(fmt.Sprintf("ATTACH DATABASE \"%s\" AS %s;", f, name)).Execute()

	return err
}

func (x *DB) Detach(name string) error {
	_, err := x.DB.NewQuery(fmt.Sprintf("DETACH DATABASE %s;", name)).Execute()

	return err
}

func (x *DB) CreateIndexIfNotExists(table, name string, cols ...string) error {
	if index := strings.Index(table, "."); index != -1 {
		table = table[index+1:]
	}
	_, err := x.DB.NewQuery(fmt.Sprintf(
		"CREATE INDEX IF NOT EXISTS %v ON %v (%v)",
		x.DB.QuoteColumnName(name),
		x.DB.QuoteTableName(table),
		x.quoteColumns(cols),
	)).Execute()

	return err
}

func (x *DB) CreateUniqueIndexIfNotExists(table, name string, cols ...string) error {
	if index := strings.Index(table, "."); index != -1 {
		table = table[index+1:]
	}
	_, err := x.DB.NewQuery(fmt.Sprintf(
		"CREATE UNIQUE INDEX IF NOT EXISTS %v ON %v (%v)",
		x.DB.QuoteColumnName(name),
		x.DB.QuoteTableName(table),
		x.quoteColumns(cols),
	)).Execute()

	return err
}

func (x *DB) quoteColumns(cols []string) string {
	s := ""
	for i, col := range cols {
		if i == 0 {
			s = x.DB.QuoteColumnName(col)
		} else {
			s += ", " + x.DB.QuoteColumnName(col)
		}
	}
	return s
}
