package ft

import (
	"fmt"
	"strings"
)

// Last iterate through iter until last element and return last element
// do not call on endless iterators
func Last[T any](iter Iter[T]) T {
	next := iter.Next()
	last := next
	for next != nil {
		last = next
		next = iter.Next()
	}
	return *last
}

// First return first element of iter
func First[T any](iter Iter[T]) T {
	return *iter.Next()
}

// Join return string of `slice` elements separated by `sep`
func Join[T any, S ~[]T](slice S, sep string) string {
	parts := make([]string, 0, len(slice))
	for _, s := range slice {
		parts = append(parts, fmt.Sprintf("%v", s))
	}
	return strings.Join(parts, sep)
}

func Contains[T comparable](iter Iter[T], elem T) bool {
	for next := iter.Next(); next != nil; next = iter.Next() {
		if elem == *next {
			return true
		}
	}
	return false
}

// Find returns first element for which predicate returns true, nil if no such element
func Find[T any](iter Iter[T], predicate func(T) bool) *T {
	for next := iter.Next(); next != nil; next = iter.Next() {
		if predicate(*next) {
			return next
		}
	}
	return nil
}
