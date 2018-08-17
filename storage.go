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
	default:
		return nil, errors.New("invalid storage engine")
	}
	return nil, errors.New("something went wrong while creating new storage")
}

// Storage is the primary interface for interacting with Payload and PayloadMetaData
type Storage interface {
	Get(key string) (Payload, PayloadMetaData, error)
	Set(key string, value Payload, meta PayloadMetaData) error
	Delete(key string) error
}

// MemoryStore  is the primary in-memory data storage and retrieval struct
type MemoryStore struct {
	sync.RWMutex
	payload  map[string]Payload
	metadata map[string]PayloadMetaData
}

// NewMemoryStore returns a pointer to a new intance ofMemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		payload:  make(map[string]Payload),
		metadata: make(map[string]PayloadMetaData),
	}
}

// Get method for MemoryStore with read locks
func (m *MemoryStore) Get(key string) (Payload, PayloadMetaData, error) {
	var err error
	m.RLock()
	payloadResult, ok := m.payload[key]
	metadataResult := m.metadata[key]
	m.RUnlock()
	if !ok {
		err = errors.New("key not found")
	}
	return payloadResult, metadataResult, err
}

// Set method for MemoryStore with read/write locks
func (m *MemoryStore) Set(key string, value Payload, meta PayloadMetaData) error {
	m.Lock()
	m.payload[key] = value
	m.metadata[key] = meta
	m.Unlock()
	return nil
}

// Delete method for MemoryStore with read/write locks
func (m *MemoryStore) Delete(key string) error {
	m.Lock()
	delete(m.payload, key)
	delete(m.metadata, key)
	m.Unlock()
	return nil
}
