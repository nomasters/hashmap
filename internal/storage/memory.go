package storage

import (
	"bytes"
	"errors"
	"sync"
	"time"
)

// MemoryStore is the primary in-memory data storage and retrieval struct. It contains
// a sync.RWMutex and an internal map of `map[string][]byte` to store state that conforms
// to the Storage interface.
type MemoryStore struct {
	sync.RWMutex
	internal map[string]memVal
}

// memVal is the value wrapper in the MemoryStore internal map and is used to
// store the timestamp of the payload to prevent replay attacks
type memVal struct {
	payload   []byte
	timestamp time.Time
}

// NewMemoryStore returns a reference to a MemoryStore with an initialized internal map
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		internal: make(map[string]memVal),
	}
}

// Get takes a key string and returns a byte slice and error. This method uses read locks.
// It returns an error if the key is not found.
func (s *MemoryStore) Get(key string) ([]byte, error) {
	s.RLock()
	v, ok := s.internal[key]
	s.RUnlock()
	if !ok {
		return []byte{}, errors.New("key not found")
	}
	return v.payload, nil
}

// Set takes a key string and byte slice value and returns an error. It uses a mutex write lock for safety.
// If an existing key value pair exists, it checks the timestamp and rejects <= timestamp submissions.
func (s *MemoryStore) Set(key string, value []byte, ttl time.Duration, timestamp time.Time) error {
	s.Lock()

	v, ok := s.internal[key]
	if ok {
		if v.timestamp.UnixNano()/1000 >= timestamp.UnixNano()/1000 {
			s.Unlock()
			return errInvalidTimestamp
		}
	}
	s.internal[key] = memVal{
		payload:   value,
		timestamp: timestamp,
	}
	s.Unlock()
	go func() {
		time.Sleep(safeTTL(ttl))
		s.deleteIfValueMatch(key, value)
	}()
	return nil
}

// Close implements the standard Close method for storage.
func (s *MemoryStore) Close() error {
	return nil
}

// deleteIfValueMatch compares values and deletes if they are they same
func (s *MemoryStore) deleteIfValueMatch(key string, value []byte) {
	s.Lock()
	if bytes.Equal(s.internal[key].payload, value) {
		delete(s.internal, key)
	}
	s.Unlock()
}
