package storage

import (
	"bytes"
	"testing"
	"time"
)

func TestMemoryStore_Get(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		expected := []byte("such_dead_much_beef")
		s := NewMemoryStore()
		defer s.Close()
		s.internal[key] = memVal{
			payload:   expected,
			timestamp: time.Now(),
		}

		actual, err := s.Get(key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(expected, actual) {
			t.Errorf("actual: %v, expected: %v", actual, expected)
		}
	})

	t.Run("without relevant entry", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		s := NewMemoryStore()
		if _, err := s.Get(key); err == nil {
			t.Fail()
		}
	})
}

func TestMemoryStore_Set(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		expected := []byte("such_dead_much_beef")
		now := time.Now()
		s := NewMemoryStore()

		if _, err := s.Get(key); err == nil {
			t.Fail()
		}

		if err := s.Set(key, expected, 1*time.Second, now); err != nil {
			t.Error(err)
		}
		actual, err := s.Get(key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(expected, actual) {
			t.Errorf("actual: %v, expected: %v", actual, expected)
		}
		if err := s.Set(key, expected, 1*time.Second, now.Add(-1*time.Nanosecond)); err == nil {
			t.Error("failed to catch stale timestamp")
		}
	})

	t.Run("deleteIfValueMatch", func(t *testing.T) {
		t.Parallel()

		key := "DEADBEEF"
		expected := []byte("such_dead_much_beef")
		s := NewMemoryStore()

		if _, err := s.Get(key); err == nil {
			t.Fail()
		}

		if err := s.Set(key, expected, 1*time.Second, time.Now()); err != nil {
			t.Error(err)
		}
		s.deleteIfValueMatch(key, expected)
		if _, err := s.Get(key); err == nil {
			t.Error("failed to delete key")
		}
	})
}
