package hashmap_test

import (
	"testing"

	"github.com/nomasters/hashmap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNaClSignEd25519_Validate(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		pubkey, err := examplePayload.SignatureBytes()
		require.NoError(t, err)

		data, err := examplePayload.DataBytes()
		require.NoError(t, err)

		s := append(pubkey, data...)
		p, err := examplePayload.PubKeyBytes()
		require.NoError(t, err)

		assert.NoError(t, hashmap.NewNaClSignEd25519(s, p).Validate())
	})

	T.Run("with invalid data", func(t *testing.T) {
		t.Parallel()
		x := []byte("nope")
		assert.Error(t, hashmap.NewNaClSignEd25519(x, x).Validate())
	})
}
