package storage

import (
	"errors"
	"time"

	"github.com/nomasters/hashmap/pkg/payload"
)

// Engine is the enum type for StorageEngine
type Engine uint8

// Enum types for Storage Engine
const (
	_ Engine = iota
	MemoryEngine
	RedisEngine
)

var (
	minTTL = payload.MaxSubmitWindow*2 + time.Second // to prevent replay attacks against SubmitWindow
	maxTTL = payload.MaxTTL
)

var (
	errInvalidTimestamp = errors.New("storage: invalid timestamp")
	errInvalidStorage   = errors.New("invalid storage engine")
)

// Getter is an interface that wraps around the standard Get method.
type Getter interface {
	Get(key string) ([]byte, error)
}

// Setter is an interface that wraps around the standard Set method. To Prevent replay attacks, the ttl
// set on key should never be less than the minimumTTL, which is 2x the payload.MaxSubmitWindow. This prevents a specific
// type of replay attack in which a short-lived TTL is set, below the SubmitWindow horizon, and a slightly older message,
// also within the SubmitWindow time horizon is replayed and accepted due to no previous key being in the system. Timestamps
// outside of this window should be automatically forward-looking, but a proper Set method should always check that the submitted
// timestamp is greater than the existing one, if an existing one is found.
type Setter interface {
	Set(key string, value []byte, ttl time.Duration, timestamp time.Time) error
}

// Closer is an interface that wraps around the standard Close method.
type Closer interface {
	Close() error
}

// GetSetCloser is the interface that groups the basic Get, Set and Close methods.
type GetSetCloser interface {
	Getter
	Setter
	Closer
}

// options is us to store Storage related Options
type options struct {
	engine Engine
	redis  []RedisOption
}

// Option is used for special Settings in Storage
type Option func(*options)

// parseOptions takes a arbitrary number of Option funcs and returns a options struct
func parseOptions(opts ...Option) options {
	var o options
	for _, option := range opts {
		option(&o)
	}
	return o
}

// New is a helper function that takes an arbitrary number of options and returns a GetSetCloser interface
func New(opts ...Option) (GetSetCloser, error) {
	o := parseOptions(opts...)
	switch o.engine {
	case MemoryEngine:
		return NewMemoryStore(), nil
	case RedisEngine:
		return NewRedisStore(o.redis...), nil
	default:
		return nil, errInvalidStorage
	}
}

// WithEngine takes an Engine and returns an Option.
func WithEngine(e Engine) Option {
	return func(o *options) {
		o.engine = e
	}
}

// WithRedisOptions takes an arbitrary number of RedisOption and returns a Option
func WithRedisOptions(opts ...RedisOption) Option {
	return func(o *options) {
		o.redis = opts
	}
}

// safeTTL ensures that a submitted TTL is no less than 2x the SubmitWindow duration, to prevent replay attacks
// it also ensures that the maxTTL is no greater than allowed.
func safeTTL(ttl time.Duration) time.Duration {
	if ttl <= minTTL {
		return minTTL
	}
	if ttl > maxTTL {
		return maxTTL
	}
	return ttl
}
