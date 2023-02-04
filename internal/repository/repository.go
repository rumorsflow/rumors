package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
)

var (
	ErrEntityNotFound       = errs.New("could not find entity")
	ErrDuplicateKey         = errs.New("entity duplicate key")
	ErrMissingEntityID      = errs.New("missing entity ID")
	ErrMissingCursor        = errs.New("missing iterator cursor")
	ErrMissingEntityFactory = errs.New("missing entity factory")
)

const (
	OpNew      errs.Op = "repository: new"
	OpIter     errs.Op = "repository: iter"
	OpFind     errs.Op = "repository: find"
	OpFindIter errs.Op = "repository: find iter"
	OpFindByID errs.Op = "repository: find by ID"
	OpCount    errs.Op = "repository: count"
	OpSave     errs.Op = "repository: save"
	OpRemove   errs.Op = "repository: remove"
	OpIndexes  errs.Op = "repository: indexes"

	ErrMsgDecode    = "failed to decode entity due to error: %w"
	ErrMsgAfterFind = "failed after find callback due to error: %w"
)

type ReadRepository[T Entity] interface {
	Count(ctx context.Context, filter any) (int64, error)
	Find(ctx context.Context, criteria *Criteria) ([]T, error)
	FindIter(ctx context.Context, criteria *Criteria) (Iter[T], error)
	FindByID(ctx context.Context, id uuid.UUID) (T, error)
}

type WriteRepository[T Entity] interface {
	Save(ctx context.Context, entity T) error
	Remove(ctx context.Context, id uuid.UUID) error
}

type ReadWriteRepository[T Entity] interface {
	ReadRepository[T]
	WriteRepository[T]
}

type innerReadRepo[T Entity] struct {
	inner        ReadRepository[T]
	changeFilter func(any) any
}

func NewReadRepository[T Entity](inner ReadRepository[T], changeFilter func(any) any) ReadRepository[T] {
	return &innerReadRepo[T]{inner: inner, changeFilter: changeFilter}
}

func (r *innerReadRepo[T]) Count(ctx context.Context, filter any) (int64, error) {
	return r.inner.Count(ctx, r.changeFilter(filter))
}

func (r *innerReadRepo[T]) Find(ctx context.Context, criteria *Criteria) ([]T, error) {
	return r.inner.Find(ctx, criteria)
}

func (r *innerReadRepo[T]) FindIter(ctx context.Context, criteria *Criteria) (Iter[T], error) {
	if criteria != nil {
		criteria.Filter = r.changeFilter(criteria.Filter)
	}
	return r.inner.FindIter(ctx, criteria)
}

func (r *innerReadRepo[T]) FindByID(ctx context.Context, id uuid.UUID) (T, error) {
	return r.inner.FindByID(ctx, id)
}
