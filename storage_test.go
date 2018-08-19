package hashmap_test

import (
	"testing"

	"github.com/nomasters/hashmap"

	"github.com/stretchr/testify/assert"
)

func TestEngine_String(T *testing.T) {
	T.Parallel()

	for engine, expected := range map[hashmap.Engine]string{
		hashmap.MemoryStorage: "Memory",
		hashmap.RedisStorage:  "Redis",
	} {
		assert.Equal(T, expected, engine.String())
	}
}

func TestNewStorage(T *testing.T) {
	T.Parallel()

	T.Run("memory storage", func(t *testing.T) {
		t.Parallel()

		_, err := hashmap.NewStorage(hashmap.MemoryStorage, hashmap.StorageOptions{})
		assert.NoError(t, err)
	})

	T.Run("redis storage", func(t *testing.T) {
		t.Parallel()

		_, err := hashmap.NewStorage(hashmap.RedisStorage, hashmap.StorageOptions{})
		assert.NoError(t, err)
	})

	T.Run("invalid storage engine", func(t *testing.T) {
		_, err := hashmap.NewStorage(666, hashmap.StorageOptions{})
		assert.Error(t, err)
	})
}
