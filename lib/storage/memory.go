package storage

import (
	"errors"
	"sync"
)

// MemoryStore is the primary in-memory data storage and retrieval struct. It contains
// a sync.RWMutex and an internal map of `map[string][]byte` to store state that conforms
// to the Storage interface.
type MemoryStore struct {
	sync.RWMutex
	internal map[string][]byte
}

// NewMemoryStore returns a reference to a MemoryStore with an initialized internal map
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		internal: make(map[string][]byte),
	}
}

// Get takes a key string and returns a byte slice and error. This method uses read locks.
// It returns an error if the key is not found.
func (s *MemoryStore) Get(key string) (v []byte, err error) {
	s.RLock()
	v, ok := s.internal[key]
	s.RUnlock()
	if !ok {
		err = errors.New("key not found")
	}
	return
}

// Set takes a key string and byte slice value and returns an error. It uses a mutex write lock for safety.
// It always returns nil.
func (s *MemoryStore) Set(key string, value []byte, options ...Option) error {
	s.Lock()
	s.internal[key] = value
	s.Unlock()
	return nil
}

// Delete takes a key string and returns an error.  It uses a mutex write lock for safety.
// It always returns nil.
func (s *MemoryStore) Delete(key string) error {
	s.Lock()
	delete(s.internal, key)
	s.Unlock()
	return nil
}
