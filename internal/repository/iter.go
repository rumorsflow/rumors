package repository

import (
	"context"
	"fmt"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
)

var _ Iter[Entity] = (*Iterator[Entity])(nil)

type Iter[T Entity] interface {
	Next(ctx context.Context) bool
	Entity() T
	Close(ctx context.Context) error
}

type Cursor interface {
	Next(ctx context.Context) bool
	Decode(val any) error
	Close(ctx context.Context) error
	Err() error
}

type Iterator[T Entity] struct {
	Cursor    Cursor
	Factory   EntityFactory[T]
	AfterFind func(entity T) error
	entity    T
	decodeErr error
}

func (i *Iterator[T]) Next(ctx context.Context) bool {
	if !i.valid() || !i.Cursor.Next(ctx) {
		return false
	}

	i.entity = i.Factory.NewEntity()

	if i.decodeErr = i.Cursor.Decode(i.entity); i.decodeErr != nil {
		i.decodeErr = fmt.Errorf(ErrMsgDecode, i.decodeErr)
		return false
	}

	if i.AfterFind != nil {
		if i.decodeErr = i.AfterFind(i.entity); i.decodeErr != nil {
			i.decodeErr = errs.E(i.entity.EntityID(), fmt.Errorf(ErrMsgAfterFind, i.decodeErr))
			return false
		}
	}

	return true
}

func (i *Iterator[T]) Entity() T {
	return i.entity
}

func (i *Iterator[T]) Close(ctx context.Context) (err error) {
	_ = i.valid()

	if i.Cursor != nil {
		err = i.Cursor.Close(ctx)
	}

	if i.decodeErr != nil || err != nil {
		err = errs.E(OpIter, i.decodeErr, err)
	}

	return
}

func (i *Iterator[T]) valid() bool {
	if i.decodeErr != nil {
		return false
	}

	if i.Cursor == nil {
		i.decodeErr = errs.E(OpFind, ErrMissingCursor)
		return false
	}

	if i.Factory == nil {
		i.decodeErr = errs.E(OpFind, ErrMissingEntityFactory)
		return false
	}

	return true
}
