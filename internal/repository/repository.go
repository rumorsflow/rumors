package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrEntityNotFound       = errors.New("could not find entity")
	ErrDuplicateKey         = errors.New("entity duplicate key")
	ErrMissingEntityID      = errors.New("missing entity ID")
	ErrMissingCursor        = errors.New("missing iterator cursor")
	ErrMissingEntityFactory = errors.New("missing entity factory")
)

const (
	OpNew      = "repository: new ->"
	OpIter     = "repository: iter ->"
	OpFind     = "repository: find ->"
	OpFindIter = "repository: find iter ->"
	OpFindByID = "repository: find by ID ->"
	OpCount    = "repository: count ->"
	OpSave     = "repository: save ->"
	OpRemove   = "repository: remove ->"
	OpIndexes  = "repository: indexes ->"

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
