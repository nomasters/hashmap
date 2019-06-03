package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	timeHorizon    = 5 * time.Second
	maxPayloadSize = 50000
)

func main() {
	m := NewMemoryStore()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for x := 0; x < 100000; x++ {
		k := fmt.Sprintf("%v", x)
		v := make([]byte, maxPayloadSize)
		d := time.Duration(r1.Intn(5)) * time.Second
		fmt.Printf("setting %v with ttl: %v\n", k, d)
		m.Set(k, v, d)
	}
	time.Sleep(6 * time.Second)
}

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
func (s *MemoryStore) Set(key string, value []byte, ttl time.Duration) error {
	s.Lock()
	s.internal[key] = value
	s.Unlock()
	go func() {
		time.Sleep(ttl + timeHorizon)
		s.DeleteIfSame(key, value)
	}()
	return nil
}

// DeleteIfSame compares values and deletes if they are they same
func (s *MemoryStore) DeleteIfSame(key string, value []byte) error {
	s.Lock()
	if bytes.Equal(s.internal[key], value) {
		delete(s.internal, key)
		fmt.Printf("deleted %v with TTL\n", key)
	}
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
