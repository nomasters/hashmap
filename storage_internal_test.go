package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore_Get(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleKey := "example"
		expected := PayloadWithMetadata{Payload: *examplePayload}
		ms := NewMemoryStore()
		ms.internal = map[string]PayloadWithMetadata{exampleKey: expected}

		actual, err := ms.Get(exampleKey)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("without relevant entry", func(t *testing.T) {
		t.Parallel()

		exampleKey := "example"
		ms := NewMemoryStore()

		_, err := ms.Get(exampleKey)

		assert.Error(t, err)
	})
}

func TestMemoryStore_Set(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleKey := "example"
		expected := PayloadWithMetadata{Payload: *examplePayload}
		ms := NewMemoryStore()

		_, err := ms.Get(exampleKey)
		assert.Error(t, err)

		ms.Set(exampleKey, expected)

		actual, err := ms.Get(exampleKey)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestMemoryStore_Delete(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleKey := "example"
		ms := NewMemoryStore()
		ms.Set(exampleKey, PayloadWithMetadata{Payload: *examplePayload})

		_, err := ms.Get(exampleKey)
		assert.NoError(t, err)

		ms.Delete(exampleKey)

		_, err = ms.Get(exampleKey)
		assert.Error(t, err)
	})
}
