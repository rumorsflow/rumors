package db

import (
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"sync"
)

type resolver[T repository.Entity] struct {
	fn   func()
	r    repository.ReadWriteRepository[T]
	err  error
	once sync.Once
}

func newResolver[T repository.Entity](fn func() (repository.ReadWriteRepository[T], error)) *resolver[T] {
	r := &resolver[T]{}
	r.fn = func() {
		r.r, r.err = fn()
	}
	return r
}

func (r *resolver[T]) Resolve() (repository.ReadWriteRepository[T], error) {
	r.once.Do(r.fn)
	return r.r, r.err
}
