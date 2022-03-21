package gsync

import "sync"

type Mutex[T any] struct {
	mutex *sync.Mutex
	val   *T
}

func NewMutex[T any](t T) *Mutex[T] {
	m := &Mutex[T]{
		val:   &t,
		mutex: &sync.Mutex{},
	}
	return m
}

func (m *Mutex[T]) Lock() *T {
	m.mutex.Lock()
	return m.val
}

func (m *Mutex[_]) Unlock() {
	m.mutex.Unlock()
}

func (m *Mutex[T]) TryLock() (*T, bool) {
	ok := m.mutex.TryLock()
	return m.val, ok
}
