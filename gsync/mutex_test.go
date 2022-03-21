package gsync_test

import (
	"gtools/gsync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Some[T any] struct {
	protected *gsync.Mutex[T]
}

func (s *Some[T]) GetProtectedData() *gsync.Mutex[T] {
	return s.protected
}

func TestMutex_Simple(t *testing.T) {
	s := Some[int]{
		protected: gsync.NewMutex(123),
	}

	lock := s.GetProtectedData()

	val := lock.Lock()
	assert.Equal(t, 123, *val)
	v, ok := lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()

	val = lock.Lock()
	*val = 42
	v, ok = lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()
	assert.Equal(t, 42, *s.protected.Lock())
	s.protected.Unlock()
}

func TestMutex_Ptr(t *testing.T) {
	a := 123
	s := Some[*int]{
		protected: gsync.NewMutex(&a),
	}

	lock := s.GetProtectedData()

	val := lock.Lock()
	assert.Equal(t, 123, **val)
	v, ok := lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()

	val = lock.Lock()
	**val = 42
	v, ok = lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()
	assert.Equal(t, 42, **s.protected.Lock())
	s.protected.Unlock()
}

func TestMutex_Nil(t *testing.T) {
	a := 42
	s := Some[*int]{
		protected: gsync.NewMutex[*int](nil),
	}

	lock := s.GetProtectedData()

	val := lock.Lock()
	assert.Nil(t, *val)
	v, ok := lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()

	val = lock.Lock()
	*val = &a
	v, ok = lock.TryLock()
	assert.False(t, ok)
	assert.Equal(t, *val, *v)
	lock.Unlock()
	assert.Equal(t, 42, **s.protected.Lock())
	s.protected.Unlock()
}
