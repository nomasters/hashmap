package storage

import (
	"bytes"
	"encoding/base64"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
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

	t.Run("DialErr", func(t *testing.T) {
		s := NewRedisStore(WithRedisEndpoint(":58555"))
		if err := s.Delete("dialFail"); err == nil {
			t.Error(err)
		}
	})

	t.Run("AuthErr", func(t *testing.T) {
		s := NewRedisStore(WithRedisEndpoint(r.Addr()), WithRedisAuth("wrong_auth"))
		if err := s.Delete("authFail"); err == nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		key := "get1"
		expected := []byte("exp1")
		encoded := base64.StdEncoding.EncodeToString(expected)
		r.Set(key, encoded)

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
			}
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("actual: %v, expected: %v description: %v\n", actual, test.expected, test.description)
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		tests := []struct {
			key         string
			expected    []byte
			ttl         Option
			fastForward time.Duration
			description string
		}{
			{
				key:         "set1",
				expected:    []byte("exp1"),
				ttl:         WithTTL(0 * time.Second),
				fastForward: (defaultRedisTTL + 1) * time.Second,
				description: "should set expected value",
			},
			{
				key:         "set2",
				expected:    []byte("exp2"),
				ttl:         WithTTL((maxRedisTTL + 30) * time.Second),
				fastForward: (maxRedisTTL + 1) * time.Second,
				description: "should conform to maxRedisTTL",
			},
			{
				key:         "set3",
				expected:    []byte("exp3"),
				ttl:         WithTTL(1 * time.Second),
				fastForward: 2 * time.Second,
				description: "should conform to short TTL",
			},
		}

		for _, test := range tests {
			if r.Exists(test.key) {
				t.Errorf("key: %v should not exist", test.key)
				continue
			}
			err := s.Set(test.key, test.expected, test.ttl)
			if err != nil {
				t.Error(test.description, err)
				continue
			}
			encoded := base64.StdEncoding.EncodeToString(test.expected)
			r.CheckGet(t, test.key, encoded)
			r.FastForward(test.fastForward)
			if _, err := r.Get(test.key); err == nil {
				t.Errorf("key: %v should not exist after ttl", test.key)
			}
		}
	})
	t.Run("Delete", func(t *testing.T) {
		tests := []struct {
			key         string
			value       string
			exists      bool
			description string
		}{
			{
				key:         "del1",
				value:       "exp1",
				exists:      true,
				description: "delete an existing key",
			},
			{
				key:         "del1",
				exists:      false,
				description: "delete an nonexisting key",
			},
		}

		for _, test := range tests {
			if test.exists {
				r.Set(test.key, "exp1")
				s.Delete(test.key)
				if r.Exists(test.key) {
					t.Error(test.description, test.key)
				}
				continue
			}
			if r.Exists(test.key) {
				t.Error(test.description, test.key)
				return
			}
			if err := s.Delete(test.key); err != nil {
				t.Error(test.description)
			}

		}
	})
}
