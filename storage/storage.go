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

// Options is us to store Storage related Options
type Context struct {
	engine Engine
	redis  RedisOptions
	ttl    time.Duration
}

// Option is used for special Settings in Storage
type Option func(*Context)

// parseOptions takes a arbitrary number of Option funcs and returns an Options struct
func parseOptions(options ...Option) Context {
	var c Context
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
	return func(c *Context) {
		c.ttl = d
	}
}

// WithEngine takes an Engine and returns an Option.
func WithEngine(e Engine) Option {
	return func(c *Context) {
		c.engine = e
	}
}
