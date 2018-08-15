package hashmap

import (
	"sync"
)

// HashCache  is the primary in-memory data storage and retrieval struct
type HashCache struct {
	sync.RWMutex
	internal map[string]Payload
}

// NewHashCache returns a pointer to a new intance of HashCache
func NewHashCache() *HashCache {
	return &HashCache{
		internal: make(map[string]Payload),
	}
}

// Get method for HashCache with read locks
func (hc *HashCache) Get(key string) (Payload, bool) {
	hc.RLock()
	result, ok := hc.internal[key]
	hc.RUnlock()
	return result, ok
}

// Set method for HashCache with read/write locks
func (hc *HashCache) Set(key string, value Payload) {
	hc.Lock()
	hc.internal[key] = value
	hc.Unlock()
}

// Delete method for HashCache with read/write locks
func (hc *HashCache) Delete(key string) {
	hc.Lock()
	delete(hc.internal, key)
	hc.Unlock()
}
