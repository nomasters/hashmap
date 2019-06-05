package storage

import (
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	t.Run("memory storage", func(t *testing.T) {
		t.Parallel()

		if _, err := New(WithEngine(MemoryEngine)); err != nil {
			t.Error(err)
		}
	})

	t.Run("redis storage", func(t *testing.T) {
		t.Parallel()

		if _, err := New(
			WithEngine(RedisEngine),
			WithRedisOptions(WithRedisEndpoint("")),
		); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid storage engine", func(t *testing.T) {
		t.Parallel()

		if _, err := New(WithEngine(0)); err == nil {
			t.Error("failed to error with invalid engine")
		}
	})

	t.Run("safeTTL", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			timestamp time.Duration
			expected  time.Duration
			message   string
		}{
			{
				timestamp: time.Second,
				expected:  minTTL,
				message:   "failed to enforce minTTL",
			},
			{
				timestamp: maxTTL + time.Hour,
				expected:  maxTTL,
				message:   "failed to enforce maxTTL",
			},
			{
				timestamp: minTTL + time.Minute,
				expected:  minTTL + time.Minute,
				message:   "failed to allow valid",
			},
		}

		for _, test := range tests {
			t.Log(test.timestamp)
			t.Log(safeTTL(test.timestamp))
			t.Log(test.expected)
			if safeTTL(test.timestamp) != test.expected {
				t.Error(test.message)
			}
		}
	})
}
