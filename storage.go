package hashmap

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Engine int

type StorageOptions struct {
	Address string
}

const (
	MemoryStorage Engine = iota
	RedisStorage
)

const MetadataPrefix = "meta-"

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

		return NewRedisStore(opts), nil
	default:
		return nil, errors.New("invalid storage engine")
	}
	return nil, errors.New("something went wrong while creating new storage")
}

// Storage is the primary interface for interacting with Payload and PayloadMetaData
type Storage interface {
	Get(key string) (PayloadWithMetadata, error)
	Set(key string, value PayloadWithMetadata) error
	Delete(key string) error
}

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

// RedisStore is a struct with methods that conforms to the Storage Interface
type RedisStore struct {
	pool *redis.Pool
}

func NewRedisStore(opts StorageOptions) *RedisStore {
	return &RedisStore{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", opts.Address) },
		},
	}
}

// Get method for RedisStore
func (r *RedisStore) Get(key string) (PayloadWithMetadata, error) {
	var pwm PayloadWithMetadata
	c := r.pool.Get()
	defer c.Close()

	response, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return pwm, err
	}

	log.Println(response)

	mp := []byte(response["payload"])
	p := Payload{}
	if err := json.Unmarshal(mp, &p); err != nil {
		return pwm, err
	}

	pwm.Payload = p

	for k, v := range response {
		if strings.HasPrefix(k, MetadataPrefix) {
			pwm.Metadata[k[len(MetadataPrefix):]] = v
		}
	}

	return pwm, nil
}

// Set method of RedisStore
func (r *RedisStore) Set(key string, value PayloadWithMetadata) error {

	data, err := value.Payload.GetData()
	if err != nil {
		return err
	}

	ttl := data.TTL
	mp, err := json.Marshal(value.Payload)
	if err != nil {
		return err
	}
	c := r.pool.Get()
	defer c.Close()
	log.Println("made it here")
	// atomic blog for writing all hash values and then setting TTL
	c.Send("MULTI")
	c.Send("HSET", key, "payload", string(mp))

	for k, v := range value.Metadata {
		mk := MetadataPrefix + k
		c.Send("HSET", key, mk, v)
	}

	c.Send("EXPIRE", key, ttl)
	if _, err := c.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

// Delete method for RedisStore
func (r *RedisStore) Delete(key string) error {
	c := r.pool.Get()
	defer c.Close()

	if _, err := c.Do("DEL", key); err != nil {
		return err
	}

	return nil
}
