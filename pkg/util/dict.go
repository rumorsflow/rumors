package util

import (
	"sync"
	"sync/atomic"
)

type Dict[T any] struct {
	data *sync.Map
	size atomic.Int32
}

func NewDict[T any]() *Dict[T] {
	return &Dict[T]{data: &sync.Map{}}
}

func (d *Dict[T]) Del(key string) (value T) {
	if item, ok := d.data.LoadAndDelete(key); ok {
		d.size.Add(-1)
		return item.(T)
	}
	return
}

func (d *Dict[T]) Set(key string, item T) {
	d.data.Store(key, item)
	d.size.Add(1)
}

func (d *Dict[T]) Get(key string) (value T) {
	if v, ok := d.data.Load(key); ok {
		value = v.(T)
	}
	return
}

func (d *Dict[T]) Len() int {
	return int(d.size.Load())
}

func (d *Dict[T]) Has(key string) bool {
	_, ok := d.data.Load(key)
	return ok
}

func (d *Dict[T]) Loop(fn func(string, T) bool) {
	d.data.Range(func(key, value any) bool {
		return fn(key.(string), value.(T))
	})
}
