package hashmap

import (
	"errors"
	"sync"
)

type Engine int
type StorageOptions map[string]string

const (
	MemoryStorage Engine = iota
	RedisStorage
)

// String is used to pretty print storage engine constants
func (s Engine) String() string {
	names := []string{
		"Memory",
		"Redis",
	}
	return names[s]
}

// NewStorage is a helper function used for configuring supported storage engines
func NewStorage(e Engine, opts StorageOptions) (Storage, error) {
	switch e {
	case MemoryStorage:
		return NewMemoryStore(), nil
	case RedisStorage:
		return nil, nil
	default:
		return nil, errors.New("invalid storage engine")
	}
}

// Storage is the primary interface for interacting with Payload and PayloadMetaData
type Storage interface {
	Get(key string) (PayloadWithMetadata, error)
	Set(key string, value PayloadWithMetadata) error
	Delete(key string) error
}

var _ Storage = (*MemoryStore)(nil)

// MemoryStore  is the primary in-memory data storage and retrieval struct
type MemoryStore struct {
	sync.RWMutex
	internal map[string]PayloadWithMetadata
}

// NewMemoryStore returns a pointer to a new intance ofMemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		internal: make(map[string]PayloadWithMetadata),
	}
}

// Get method for MemoryStore with read locks
func (m *MemoryStore) Get(key string) (PayloadWithMetadata, error) {
	var err error
	m.RLock()
	v, ok := m.internal[key]
	m.RUnlock()
	if !ok {
		err = errors.New("key not found")
	}
	return v, err
}

// Set method for MemoryStore with read/write locks
func (m *MemoryStore) Set(key string, value PayloadWithMetadata) error {
	m.Lock()
	m.internal[key] = value
	m.Unlock()
	return nil
}

// Delete method for MemoryStore with read/write locks
func (m *MemoryStore) Delete(key string) error {
	m.Lock()
	delete(m.internal, key)
	m.Unlock()
	return nil
}
