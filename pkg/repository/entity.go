package repository

import (
	"github.com/google/uuid"
	"reflect"
)

type Entity interface {
	EntityID() uuid.UUID
}

type EntityFactory[T Entity] interface {
	NewEntity() T
}

type EntityFactoryFunc[T Entity] func() T

func (f EntityFactoryFunc[T]) NewEntity() T {
	return f()
}

func Factory[T Entity]() EntityFactory[T] {
	return EntityFactoryFunc[T](func() (value T) {
		v := reflect.ValueOf(value)
		value = reflect.New(v.Type().Elem()).Interface().(T)
		return
	})
}
