package hashmap

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/sign"
)

const (
	exampleValidPayload = `
		{
			"data": "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			"sig": "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			"pubkey": "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo="
		}
	`
)

var (
	examplePayload = &Payload{
		Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
		Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
		PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
	}
)

// Key returns a randomly generated ed25519 private key in bytes
func buildInvalidKey() string {
	_, privKey, _ := sign.GenerateKey(rand.Reader)
	return base64.StdEncoding.EncodeToString(privKey[31:])
}

type errorReader struct{}

func (er errorReader) Read(b []byte) (int, error) {
	return 0, errors.New("arbitrary")
}

func TestNewPayloadFromReader(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		actual, err := NewPayloadFromReader(strings.NewReader(exampleValidPayload))
		assert.NoError(t, err)
		assert.Equal(t, examplePayload, actual)
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()

		shouldBeNil, err := NewPayloadFromReader(&errorReader{})
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with invalid JSON", func(t *testing.T) {
		t.Parallel()

		exampleReader := strings.NewReader("not json lol")
		shouldBeNil, err := NewPayloadFromReader(exampleReader)
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with payload validation error", func(t *testing.T) {
		t.Parallel()
		examplePayload := `
		{
			"data": "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			"sig": "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			"pubkey": "this is bad lol"
		}
		`
		exampleReader := strings.NewReader(examplePayload)

		shouldBeNil, err := NewPayloadFromReader(exampleReader)
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestNewValidator(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		v, err := examplePayload.NewValidator()

		assert.NoError(t, v.Validate())
		assert.NoError(t, err)
	})

	T.Run("with pubkey failure", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "nope lol",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
	T.Run("with signature failure", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "nope lol",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
	T.Run("with data failure", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
	T.Run("with GetData failure", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "8J+ZifCfmYjwn5mK", // base64 encoded monkeys
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
	T.Run("with invalid signature method", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6InNvbWV0aGluJyBlbHNlIiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}

		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
	T.Run("with too long a public key", func(t *testing.T) {
		t.Parallel()

		example := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: buildInvalidKey(),
		}

		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestPayload_PubKeyBytes(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		actual, err := examplePayload.PubKeyBytes()

		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
	})
}

func TestPayload_Verify(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, examplePayload.Verify())
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()
		badPayload := `
		{
			"data": "blarg",
			"sig": "nature",
			"pubkey": "this is bad lol"
		}
		`
		p := &Payload{}
		json.Unmarshal([]byte(badPayload), &p)

		assert.Error(t, p.Verify())
	})
}

func TestPayload_GetData(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		actual, err := examplePayload.DataBytes()
		assert.Empty(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error decoding data", func(t *testing.T) {
		example := &Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		actual, err := example.GetData()
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
