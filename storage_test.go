package hashmap_test

import (
	"testing"

	"github.com/nomasters/hashmap"

	"github.com/stretchr/testify/assert"
)

func TestEngine_String(T *testing.T) {
	T.Parallel()

	for engine, expected := range map[hashmap.Engine]string{
		hashmap.MemoryStorage: "memory",
		hashmap.RedisStorage:  "redis",
	} {
		assert.Equal(T, expected, engine.String())
	}
}

func TestNewStorage(T *testing.T) {
	T.Parallel()

	T.Run("memory storage", func(t *testing.T) {
		t.Parallel()

		_, err := hashmap.NewStorage(hashmap.StorageOptions{Engine: hashmap.MemoryStorage})
		assert.NoError(t, err)
	})

	T.Run("redis storage", func(t *testing.T) {
		t.Parallel()

		_, err := hashmap.NewStorage(hashmap.StorageOptions{Engine: hashmap.RedisStorage})
		assert.NoError(t, err)
	})

	T.Run("invalid storage engine", func(t *testing.T) {
		_, err := hashmap.NewStorage(hashmap.StorageOptions{Engine: 3})
		assert.Error(t, err)
	})
}
