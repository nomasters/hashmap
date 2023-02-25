package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestRedisStore(t *testing.T) {
	t.Parallel()
	r, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer r.Close()

	auth := "super_secret_redis_auth"
	r.RequireAuth(auth)

	s := NewRedisStore(
		WithRedisEndpoint(r.Addr()),
		WithRedisDialTLSSkipVerify(false),
		WithRedisIdleTimeout(15*time.Second),
		WithRedisMaxActive(20),
		WithRedisMaxConnLifetime(20),
		WithRedisMaxIdle(20),
		WithRedisAuth(auth),
		WithRedisWait(true),
		WithRedisTLS(false),
	)
	defer s.Close()

	t.Run("DialErr", func(t *testing.T) {
		s := NewRedisStore(WithRedisEndpoint(":58555"))
		if _, err := s.Get("dialFail"); err == nil {
			t.Error("failed to catch dial error on get")
		}
		if err := s.Set("dialFail", []byte{}, 0, time.Now()); err == nil {
			t.Error("failed to catch dial error on set")
		}
	})

	t.Run("AuthErr", func(t *testing.T) {
		s := NewRedisStore(WithRedisEndpoint(r.Addr()), WithRedisAuth("wrong_auth"))
		if _, err := s.Get("authFail"); err == nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		key := "get1"
		expected := []byte("exp1")
		timestamp := time.Now()
		v := redisVal{
			Payload:   base64.StdEncoding.EncodeToString(expected),
			Timestamp: timestamp.UnixNano(),
		}
		enc, err := json.Marshal(v)
		if err != nil {
			t.Error(err)
		}

		r.Set(key, string(enc))

		tests := []struct {
			key         string
			expected    []byte
			shouldErr   bool
			description string
		}{
			{
				key:         key,
				expected:    expected,
				description: "should get expected value",
			},
			{
				key:         "DNE",
				shouldErr:   true,
				description: "should error an nil key",
			},
		}

		for _, test := range tests {
			actual, err := s.Get(test.key)
			if test.shouldErr {
				if err == nil {
					t.Error(test.description)
				}
				continue
			}
			if err != nil {
				t.Error(test.description, err)
				continue
			}
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("actual: %v, expected: %s description: %v\n", actual, test.expected, test.description)
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		now := time.Now()

		tests := []struct {
			key         string
			expected    []byte
			ttl         time.Duration
			timestamp   time.Time
			fastForward time.Duration
			description string
			shouldErr   bool
		}{
			{
				key:         "set1",
				expected:    []byte("exp1"),
				ttl:         0 * time.Second,
				timestamp:   now,
				fastForward: safeTTL(0*time.Second) + time.Second,
				description: "should conform to minTTL",
			},
			{
				key:         "set2",
				expected:    []byte("exp2"),
				ttl:         maxTTL + 30*time.Second,
				timestamp:   now,
				fastForward: maxTTL + time.Second,
				description: "should conform to maxTTL",
			},
			{
				key:         "set3",
				expected:    []byte("exp3"),
				ttl:         minTTL + time.Second,
				timestamp:   now,
				fastForward: minTTL + 2*time.Second,
				description: "should conform to standard TTL",
			},
		}

		for _, test := range tests {
			if r.Exists(test.key) {
				t.Errorf("key: %v should not exist", test.key)
				continue
			}
			err := s.Set(test.key, test.expected, test.ttl, test.timestamp)
			if test.shouldErr {
				if err == nil {
					t.Error(test.description)
				}
				continue
			}
			if err != nil {
				t.Error(test.description, err)
				continue
			}
			if _, err := s.Get(test.key); err != nil {
				t.Error(test.key, "failed to get")
			}
			r.FastForward(test.fastForward)
			if _, err := s.Get(test.key); err == nil {
				t.Errorf("key: %v should not exist after ttl", test.key)
			}
		}
	})
	t.Run("CheckTimeStampReplay", func(t *testing.T) {
		k := "replayattack"
		v1 := []byte("hello, world")
		v2 := []byte("bad actor")
		ttl := time.Second
		now := time.Unix(0, 1559706661127858000)
		s.Set(k, v1, ttl, now)
		if err := s.Set(k, v2, ttl, now.Add(-time.Millisecond)); err == nil {
			t.Error("failed to catch invalid timestamp")
		}
	})
	t.Run("Malformed_Get", func(t *testing.T) {
		k := "invalidGet"
		r.Set(k, "malformed")
		if _, err := s.Get(k); err == nil {
			t.Error("failed to catch malformed value")
		}
	})
}
