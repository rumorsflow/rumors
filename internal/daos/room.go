package daos

import (
	"context"
	"database/sql"
	"github.com/iagapie/rumors/internal/models"
	"github.com/iagapie/rumors/pkg/slice"
	"github.com/pocketbase/dbx"
)

var roomModel models.Room

type FilterRooms struct {
	Ids       []int64 `json:"ids,omitempty" query:"ids[]"`
	Title     *string `json:"title,omitempty" query:"title"`
	Broadcast *bool   `json:"broadcast,omitempty" query:"broadcast"`
	Deleted   *bool   `json:"deleted,omitempty" query:"deleted"`
}

func (filter *FilterRooms) Where(q *dbx.SelectQuery) *dbx.SelectQuery {
	if len(filter.Ids) > 0 {
		q.AndWhere(dbx.In(roomModel.TableName()+".id", slice.ToAny(filter.Ids)...))
	}
	if filter.Title != nil {
		q.AndWhere(like(roomModel.TableName()+".title", filter.Title))
	}
	if filter.Broadcast != nil {
		q.AndWhere(dbx.HashExp{roomModel.TableName() + ".broadcast": *filter.Broadcast})
	}
	if filter.Deleted != nil {
		q.AndWhere(dbx.HashExp{roomModel.TableName() + ".deleted": *filter.Deleted})
	}
	return q
}

func (dao *Dao) FindRoomById(ctx context.Context, id int64) (*models.Room, error) {
	data, err := dao.FindRooms(ctx, FilterRooms{Ids: []int64{id}}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return &data[0], nil
	}
	return nil, sql.ErrNoRows
}

func (dao *Dao) FindRooms(ctx context.Context, filter FilterRooms, index uint64, size uint32) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var data []models.Room

	err := filter.Where(dao.db.Select().From(roomModel.TableName())).
		Offset(int64(index)).
		Limit(int64(size)).
		WithContext(ctx).
		All(&data)

	return data, err
}

func (dao *Dao) CountRooms(ctx context.Context, filter FilterRooms) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var total uint64
	err := filter.Where(dao.db.Select("COUNT(*)").From(roomModel.TableName())).WithContext(ctx).Row(&total)

	return total, err
}
