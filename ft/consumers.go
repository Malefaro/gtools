package ft

import (
	"context"
	"reflect"
)

// CollectR same as CollectInto but this func returns value of type R (instead of filling provided arguments)
// generic types must be specified for this func during call
// type R must implement FromIter interface
// this method uses reflect to initialize nil pointers, so its less efficient than `CollectInto`
// example:
//	type SliceWrapper[T any] struct {
//		data []T
//	}
//	func (l *SliceWrapper[T]) FromIter(iter ft.Iter[T]) {
//		l.data = ft.Collect(iter)
//	}
//	iter2 := ft.SliceIter([]string{"one", "two", "three", "four"})
//	result := ft.CollectR[string, *SliceWrapper[string]](iter2)
//	// result: &main.SliceWrapper[string]{data:[]string{"one", "two", "three", "four"}}
//	fmt.Printf("result: %#v\n", result)
func CollectR[T any, R FromIter[T]](iter Iter[T]) R {
	var r R
	rv := reflect.ValueOf(&r)
	if rv.Elem().Kind() == reflect.Ptr && rv.Elem().IsNil() {
		// initialize nil pointer
		va := rv.Elem()
		v := reflect.New(va.Type().Elem())
		va.Set(v)
	}
	r.FromIter(iter)
	return r
}

// CollectInto consumes the iterator and fill provided arguments `r`
// type R must implement FromIter interface
func CollectInto[T any, R FromIter[T]](iter Iter[T], r R) {
	r.FromIter(iter)
}

// Collect consumes iter and return slice of iterator elements
func Collect[T any](iter Iter[T]) []T {
	result := make([]T, 0)
	next, ok := iter.Next()
	for ok {
		result = append(result, next)
		next, ok = iter.Next()
	}
	return result
}

// Any consumes iter and returns true if any element of iter returns true on predicate func call on it
func Any[T any](iter Iter[T], predicate func(T) bool) bool {
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		if predicate(next) {
			return true
		}
	}
	return false
}

// All consumes iter and returns true if all elements returns true on predicate func call on it
func All[T any](iter Iter[T], predicate func(T) bool) bool {
	flag := false
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		if !predicate(next) {
			return false
		}
		flag = true
	}
	if !flag {
		// if empty iterator
		return false
	}
	return true
}

// Reduce consumes iter, calls f func on every element of iterator
// f takes 2 arguments. Produced value after previous iteration and current value
// Reduce has optional argument `initial` to set initial value.
// beware of pointer types O because initial value for pointers is nil (you can check it in `f` func or set `initial`)
func Reduce[T any, O any](iter Iter[T], f func(O, T) O, initial ...O) O {
	var result O
	if len(initial) > 0 {
		result = initial[0]
	}
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		result = f(result, next)
	}
	return result
}

// Count cunsumes iter and returns count of iter elements
// if optional argument `predicate` is provided count only if predicate returns true
func Count[T any](iter Iter[T], predicate ...func(T) bool) int {
	next, ok := iter.Next()
	cnt := 0
	f := func(T) bool {
		return true
	}
	if len(predicate) > 0 {
		f = predicate[0]
	}
	for ok {
		if f(next) {
			cnt++
		}
		next, ok = iter.Next()
	}
	return cnt
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// Sum consumes iter and returns sum of elements
// work only with Numbers
// if you need sum some custom types check Reduce func
func Sum[T Number](iter Iter[T], initial ...T) T {
	next, ok := iter.Next()
	var result T
	if len(initial) > 0 {
		result = initial[0]
	}
	for ok {
		result += next
		next, ok = iter.Next()
	}
	return result
}

// ForEach consumes iter and calls func `f` on each element of iterator
func ForEach[T any](iter Iter[T], f func(T)) {
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		f(next)
	}
}

// IntoChannel converts provided iter to channel
// optional arg `ctxArg` used for stop iterations and close result channel
// context may have optional Value "size" stored in it (using `context.WithValue`)
// this value determine channel size
// maybe used for iterating through iter in `for ... := range` loop
func IntoChannel[T any](iter Iter[T], ctxArg ...context.Context) <-chan T {
	chanSize := 0
	var ctx context.Context
	if len(ctxArg) > 0 && ctxArg[0] != nil {
		ctx = ctxArg[0]
		size := ctx.Value("size")
		if size != nil {
			chanSize = size.(int)
		}
	}
	ch := make(chan T, chanSize)
	go func() {
		defer close(ch)
		for {
			next, ok := iter.Next()
			if !ok {
				return
			}
			if ctx != nil {
				select {
				case <-ctx.Done():
					return
				case ch <- next:
				}
			} else {
				ch <- next
			}
		}
	}()
	return ch
}

// Max consumes iterator and return maximum value find in it (or nil if iterator is empty)
// `less` func compare two values and return true if a < b
func Max[T any](iter Iter[T], less func(a T, b T) bool) *T {
	next, ok := iter.Next()
	if !ok {
		return nil
	}
	max := next
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		if less(max, next) {
			max = next
		}
	}
	return &max
}

// Min consumes iterator and return minimum value find in it (or nil if iterator is empty)
// `less` func compare two values and return true if a < b
func Min[T any](iter Iter[T], less func(a T, b T) bool) *T {
	next, ok := iter.Next()
	if !ok {
		return nil
	}
	min := next
	for next, ok := iter.Next(); ok; next, ok = iter.Next() {
		if less(next, min) {
			min = next
		}
	}
	return &min
}
