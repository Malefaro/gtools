package gsync_test

import (
	"gtools/gsync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeRW struct {
	protected *gsync.RWMutex[int]
}

func (s *SomeRW) GetProtectedData() *gsync.RWMutex[int] {
	return s.protected
}

func Test_RWMutexExampleLock(t *testing.T) {
	s := SomeRW{
		protected: gsync.NewRWMutex(123),
	}

	lock := s.GetProtectedData()

	val := lock.Lock()
	assert.Equal(t, 123, *val)
	lock.Unlock()

	val = lock.Lock()
	*val = 42
	lock.Unlock()
	assert.Equal(t, 42, *s.protected.Lock())
	s.protected.Unlock()
}

func Test_RWMutexExampleRLock(t *testing.T) {
	s := SomeRW{
		protected: gsync.NewRWMutex(123),
	}

	lock := s.GetProtectedData()

	val := lock.Lock()
	assert.Equal(t, 123, *val)
	lock.Unlock()

	val = lock.Lock()
	*val = 42
	lock.Unlock()
	assert.Equal(t, 42, *s.protected.Lock())
	s.protected.Unlock()
}
