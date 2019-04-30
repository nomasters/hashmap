package storage

import (
	"testing"
)

func TestNewStorage(T *testing.T) {
	T.Parallel()

	T.Run("memory storage", func(t *testing.T) {
		t.Parallel()

		if _, err := NewStorage(WithEngine(MemoryEngine)); err != nil {
			t.Error(err)
		}
	})

	T.Run("redis storage", func(t *testing.T) {
		t.Parallel()

		if _, err := NewStorage(WithEngine(RedisEngine)); err != nil {
			t.Error(err)
		}
	})

	T.Run("invalid storage engine", func(t *testing.T) {
		t.Parallel()

		if _, err := NewStorage(WithEngine(0)); err == nil {
			t.Error("failed to error with invalid engine")
		}
	})
}
