package daos

import (
	"context"
	"github.com/pocketbase/dbx"
	"strings"
	"time"
)

const timeout = 5 * time.Second

func like(col string, str *string) dbx.Expression {
	if str == nil || len(*str) == 0 {
		return nil
	}

	left, right := '%' == (*str)[0], '%' == (*str)[len(*str)-1]
	if !left && !right {
		left = true
		right = true
	}
	value := strings.ReplaceAll(*str, "%", "")

	if len(value) > 0 {
		return dbx.Like(col, value).Match(left, right)
	}

	return nil
}

type Dao struct {
	db *dbx.DB
}

func New(db *dbx.DB) *Dao {
	return &Dao{db: db}
}

func (dao *Dao) DB() *dbx.DB {
	return dao.db
}

func (dao *Dao) Insert(ctx context.Context, model any) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return dao.db.Model(model).WithContext(ctx).Insert()
}

func (dao *Dao) Update(ctx context.Context, model any) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return dao.db.Model(model).WithContext(ctx).Update()
}

func (dao *Dao) Delete(ctx context.Context, model any) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return dao.db.Model(model).WithContext(ctx).Delete()
}
