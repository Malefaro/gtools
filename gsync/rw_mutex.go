package gsync

import "sync"

type RWMutex[T any] struct {
	mutex *sync.RWMutex
	val   *T
}

func NewRWMutex[T any](t T) *RWMutex[T] {
	m := &RWMutex[T]{
		val:   &t,
		mutex: &sync.RWMutex{},
	}
	return m
}

func (m *RWMutex[T]) Lock() *T {
	m.mutex.Lock()
	return m.val
}

func (m *RWMutex[T]) Unlock() {
	m.mutex.Unlock()
}

func (m *RWMutex[T]) RLock() T {
	m.mutex.RLock()
	return *m.val
}

func (m *RWMutex[T]) RUnlock() {
	m.mutex.RUnlock()
}

func (m *RWMutex[T]) TryLock() (*T, bool) {
	ok := m.mutex.TryLock()
	return m.val, ok
}

func (m *RWMutex[T]) TryRLock() (T, bool) {
	ok := m.mutex.TryRLock()
	return *m.val, ok
}
