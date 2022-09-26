package daos

import (
	"context"
	"fmt"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/pocketbase/dbx"
)

var feedItemModel models.FeedItem

type FilterFeedItems struct {
	Ids      []int64 `json:"ids,omitempty" query:"ids[]"`
	FeedIds  []int64 `json:"feedIds,omitempty" query:"feedIds[]"`
	Author   *string `json:"author,omitempty" query:"author"`
	Category *string `json:"category,omitempty" query:"category"`
}

func (filter *FilterFeedItems) Where(q *dbx.SelectQuery) *dbx.SelectQuery {
	if len(filter.Ids) > 0 {
		q.AndWhere(dbx.In(feedItemModel.TableName()+".id", slice.ToAny(filter.Ids)...))
	}
	if len(filter.FeedIds) > 0 {
		q.AndWhere(dbx.In(feedItemModel.TableName()+".feedId", slice.ToAny(filter.FeedIds)...))
	}
	if l := like("jsonAuthor.value", filter.Author); l != nil {
		q.InnerJoin("json_each("+feedItemModel.TableName()+".authors) as jsonAuthor", l)
	}
	if l := like("jsonCategory.value", filter.Category); l != nil {
		q.InnerJoin("json_each("+feedItemModel.TableName()+".categories) as jsonCategory", l)
	}
	return q
}

func (dao *Dao) FindFeedItems(ctx context.Context, filter FilterFeedItems, index uint64, size uint32) ([]models.FeedItem, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var data []models.FeedItem

	selectCols := dao.db.QuoteColumnName(feedItemModel.Table() + ".*")
	err := filter.Where(dao.db.Select(selectCols).From(feedItemModel.TableName())).
		OrderBy(feedItemModel.TableName() + ".pubDate DESC").
		Offset(int64(index)).
		Limit(int64(size)).
		WithContext(ctx).
		All(&data)

	return data, err
}

func (dao *Dao) CountFeedItems(ctx context.Context, filter FilterFeedItems) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var total uint64
	selectCols := fmt.Sprintf("COUNT(%s)", dao.db.QuoteColumnName(feedItemModel.Table()+".id"))
	err := filter.Where(dao.db.Select(selectCols).From(feedItemModel.TableName())).WithContext(ctx).Row(&total)

	return total, err
}
