package util

import (
	"golang.org/x/exp/slices"
	"sync"
)

type List[T any] struct {
	sync.RWMutex
	data []T
}

func NewList[T any](size int) *List[T] {
	if size == 0 {
		return &List[T]{}
	}
	return &List[T]{data: make([]T, 0, size)}
}

func (l *List[T]) Len() int {
	l.RLock()
	defer l.RUnlock()

	return len(l.data)
}

func (l *List[T]) IsEmpty() bool {
	return l.Len() == 0
}

func (l *List[T]) Get(i int) T {
	l.RLock()
	defer l.RUnlock()

	return l.data[i]
}

func (l *List[T]) Shift() T {
	l.Lock()
	defer l.Unlock()

	item := l.data[0]
	l.del(0, 1)
	return item
}

func (l *List[T]) Unshift(item T) {
	l.Insert(0, item)
}

func (l *List[T]) Insert(i int, item T) {
	l.Lock()
	defer l.Unlock()

	l.data = slices.Insert(l.data, i, item)
}

func (l *List[T]) Del(i int, j int) {
	l.Lock()
	defer l.Unlock()

	l.del(i, j)
}

func (l *List[T]) Add(item T) {
	l.Lock()
	defer l.Unlock()

	l.data = append(l.data, item)
}

func (l *List[T]) del(i int, j int) {
	var item T
	l.data[i] = item
	l.data = slices.Delete(l.data, i, j)
}

func (l *List[T]) Loop(fn func(i int, item T) bool) {
	l.RLock()
	data := l.data
	l.RUnlock()

	for i, item := range data {
		if !fn(i, item) {
			break
		}
	}
}
