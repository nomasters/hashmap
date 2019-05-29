package storage

import (
	"errors"
	"time"
)

// Engine is the enum type for StorageEngine
type Engine uint8

// Enum types for Storage Engine
const (
	_ Engine = iota
	MemoryEngine
	RedisEngine
)

// Storage is the primary interface for interacting with Payload and PayloadMetaData
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, options ...Option) error
	Delete(key string) error
}

// context is us to store Storage related Options
type context struct {
	engine Engine
	redis  redisOptions
	ttl    time.Duration
}

// Option is used for special Settings in Storage
type Option func(*context)

// parseOptions takes a arbitrary number of Option funcs and returns a context struct
func parseOptions(options ...Option) context {
	var c context
	for _, option := range options {
		option(&c)
	}
	return c
}

// NewStorage is a helper function used for configuring supported storage engines
func NewStorage(options ...Option) (Storage, error) {
	o := parseOptions(options...)
	switch o.engine {
	case MemoryEngine:
		return NewMemoryStore(), nil
	case RedisEngine:
		return NewRedisStore(options...), nil
	default:
		return nil, errors.New("invalid storage engine")
	}
}

// WithTTL takes a time.Duration and returns a Option used for settings Storage related options
func WithTTL(d time.Duration) Option {
	return func(c *context) {
		c.ttl = d
	}
}

// WithEngine takes an Engine and returns an Option.
func WithEngine(e Engine) Option {
	return func(c *context) {
		c.engine = e
	}
}
