package storage

import (
	"bytes"
	"testing"
)

func TestMemoryStore_Get(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		expected := []byte("such_dead_much_beef")
		s := NewMemoryStore()
		s.internal[key] = expected

		actual, err := s.Get(key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(expected, actual) {
			t.Errorf("actual: %v, expected: %v", actual, expected)
		}
	})

	T.Run("without relevant entry", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		s := NewMemoryStore()
		if _, err := s.Get(key); err == nil {
			t.Fail()
		}
	})
}

func TestMemoryStore_Set(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		expected := []byte("such_dead_much_beef")
		s := NewMemoryStore()

		if _, err := s.Get(key); err == nil {
			t.Fail()
		}

		s.Set(key, expected)
		actual, err := s.Get(key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(expected, actual) {
			t.Errorf("actual: %v, expected: %v", actual, expected)
		}
	})
}

func TestMemoryStore_Delete(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		value := []byte("such_dead_much_beef")
		s := NewMemoryStore()
		s.Set(key, value)
		if _, err := s.Get(key); err != nil {
			t.Error(err)
		}
		s.Delete(key)
		if _, err := s.Get(key); err == nil {
			t.Error(err)
		}
	})
}
