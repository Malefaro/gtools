package ft

import (
	"reflect"
)

// Iter is the main interface for iterators
type Iter[T any] interface {
	// Next returns next iterator element. If no elements in iter return false bool flag
	Next() (T, bool)
}

// ReversibleIter interface for two direction iterators
type ReversibleIter[T any] interface {
	Iter[T]
	// Prev returns previous element of iterator
	// If it called after Next() its returns the same element as Next()
	// (because after call Next() iterator moves forward and Prev bring it back to the same element)
	Prev() (T, bool)
}

// FromIter interface used for converting iterators into structs
// Used in functions like CollectInto and CollectR
type FromIter[T any] interface {
	FromIter(iter Iter[T])
}

type sliceIter[T any] struct {
	data []T
	idx  int
}

func (si *sliceIter[T]) Next() (T, bool) {
	if si.idx > len(si.data)-1 {
		var t T
		return t, false
	}
	si.idx++
	return si.data[si.idx-1], true
}

func (si *sliceIter[T]) Prev() (T, bool) {
	if si.idx > len(si.data)-1 {
		si.idx = len(si.data)
	}
	if si.idx-1 < 0 {
		var t T
		return t, false
	}
	si.idx--
	return si.data[si.idx], true
}

// SliceIter converts provided slice into iterator
func SliceIter[T any, S ~[]T](d S) Iter[T] {
	return &sliceIter[T]{
		data: d,
		idx:  0,
	}
}

type MapPair[K comparable, V any] struct {
	Key   K
	Value V
}

type hashMapIter[K comparable, V any, R MapPair[K, V]] struct {
	data  map[K]V
	pairs chan R
}

func (mi *hashMapIter[K, V, R]) processMap() {
	for k, v := range mi.data {
		mi.pairs <- R{Key: k, Value: v}
	}
	close(mi.pairs)
}

func (mi *hashMapIter[K, V, R]) Next() (R, bool) {
	next, ok := <-mi.pairs
	if ok {
		// channel not closed
		return next, true
	} else {
		// channel closed
		var r R
		return r, false
	}
}

// MapIter returns iterator over map
// WARNING: this function runs goroutine that reads from map `m`
// To prevent goroutine leaks you must consume this iterator!
// Also you can use MapIterOverSlice to not use goroutines
func MapIter[K comparable, V any, M ~map[K]V](m M) Iter[MapPair[K, V]] {
	mi := &hashMapIter[K, V, MapPair[K, V]]{
		data:  m,
		pairs: make(chan MapPair[K, V]),
	}
	go mi.processMap()
	return mi
}

// MapIterOverSlice returns iterator over map
// unlike to MapIter this func does not create goroutine
// instead it create slice with all pairs readed from map `m`
// in case huge map sizes you should use MapIter
func MapIterOverSlice[K comparable, V any, M ~map[K]V](m M) Iter[MapPair[K, V]] {
	pairs := make([]MapPair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, MapPair[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return &sliceIter[MapPair[K, V]]{
		data: pairs,
	}
}

// copyIter make a copy(not deep) of provided iter (copy internal state of iterator) using reflect
// this func applicable to SliceIter (because it's just copy pointer to same slice and current idx)
// but maybe not applicable to some other iterators
func copyIter[T any](iter Iter[T]) Iter[T] {
	var newIter Iter[T]
	if reflect.TypeOf(iter).Kind() == reflect.Ptr {
		// Pointer:
		iterElem := reflect.ValueOf(iter).Elem()
		rIter := reflect.New(iterElem.Type())
		rIter.Elem().Set(reflect.ValueOf(iter).Elem())
		newIter = rIter.Interface().(Iter[T])
	} else {
		// Not pointer:
		return iter
	}
	return newIter
}
