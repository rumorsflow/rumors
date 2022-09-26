package daos

import (
	"context"
	"database/sql"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/pocketbase/dbx"
)

var feedModel models.Feed

type FilterFeeds struct {
	Ids     []int64 `json:"ids,omitempty" query:"ids[]"`
	By      *int64  `json:"by,omitempty" query:"by"`
	Host    *string `json:"host,omitempty" query:"host" validate:"omitempty,fqdn"`
	Link    *string `json:"link,omitempty" query:"link" validate:"omitempty,url"`
	Enabled *bool   `json:"enabled,omitempty" query:"enabled"`
}

func (filter *FilterFeeds) Where(q *dbx.SelectQuery) *dbx.SelectQuery {
	if len(filter.Ids) > 0 {
		q.AndWhere(dbx.In(feedModel.TableName()+".id", slice.ToAny(filter.Ids)...))
	}
	if filter.By != nil {
		q.AndWhere(dbx.HashExp{feedModel.TableName() + ".by": *filter.By})
	}
	if filter.Host != nil {
		q.AndWhere(like(feedModel.TableName()+".host", filter.Host))
	}
	if filter.Link != nil {
		q.AndWhere(dbx.HashExp{feedModel.TableName() + ".link": *filter.Link})
	}
	if filter.Enabled != nil {
		q.AndWhere(dbx.HashExp{feedModel.TableName() + ".enabled": *filter.Enabled})
	}
	return q
}

func (dao *Dao) FindFeedByLink(ctx context.Context, link *string) (*models.Feed, error) {
	data, err := dao.FindFeeds(ctx, FilterFeeds{Link: link}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return &data[0], nil
	}
	return nil, sql.ErrNoRows
}

func (dao *Dao) FindFeedById(ctx context.Context, id int64) (*models.Feed, error) {
	data, err := dao.FindFeeds(ctx, FilterFeeds{Ids: []int64{id}}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return &data[0], nil
	}
	return nil, sql.ErrNoRows
}

func (dao *Dao) FindFeeds(ctx context.Context, filter FilterFeeds, index uint64, size uint32) ([]models.Feed, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var data []models.Feed

	err := filter.Where(dao.db.Select().From(feedModel.TableName())).
		OrderBy("id ASC").
		Offset(int64(index)).
		Limit(int64(size)).
		WithContext(ctx).
		All(&data)

	return data, err
}

func (dao *Dao) CountFeeds(ctx context.Context, filter FilterFeeds) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var total uint64
	err := filter.Where(dao.db.Select("COUNT(*)").From(feedModel.TableName())).WithContext(ctx).Row(&total)

	return total, err
}
