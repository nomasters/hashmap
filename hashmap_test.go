package hashmap

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func buildExamplePayload() *Payload {
	p := &Payload{}
	json.Unmarshal([]byte(exampleValidPayload), &p)
	return p
}

type errorReader struct{}

func (er errorReader) Read(b []byte) (int, error) {
	return 0, errors.New("arbitrary")
}

func TestNewPayloadFromReader(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleReader := strings.NewReader(exampleValidPayload)
		expected := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		actual, err := NewPayloadFromReader(exampleReader)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()

		_, err := NewPayloadFromReader(&errorReader{})
		assert.Error(t, err)
	})

	T.Run("with invalid JSON", func(t *testing.T) {
		t.Parallel()

		exampleReader := strings.NewReader("not json lol")
		_, err := NewPayloadFromReader(exampleReader)
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

		_, err := NewPayloadFromReader(exampleReader)
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

		expected := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "nope lol",
		}
		_, err := expected.NewValidator()
		assert.Error(T, err)
	})
	T.Run("with signature failure", func(t *testing.T) {
		t.Parallel()

		expected := &Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "nope lol",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		_, err := expected.NewValidator()
		assert.Error(T, err)
	})
	T.Run("with data failure", func(t *testing.T) {
		t.Parallel()

		expected := &Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		_, err := expected.NewValidator()
		assert.Error(T, err)
	})
	T.Run("with GetData failure", func(t *testing.T) {
		t.Parallel()

		expected := &Payload{
			Data:      "8J+ZifCfmYjwn5mK", // base64 encoded monkeys
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		_, err := expected.NewValidator()
		assert.Error(T, err)
	})
}

// // NOTE: I wrote this and then wondered if I was actually testing anything.
// func TestPayload_PubKeyBytes(T *testing.T) {
// 	T.Parallel()

// 	T.Run("normal operation", func(t *testing.T) {
// 		x := buildExamplePayload()

// 		expected := []byte{0xf4, 0x76, 0x28, 0x14, 0xdd, 0xba, 0xae, 0x2, 0x7c, 0x35, 0x21, 0x6, 0x1e, 0xc8, 0x65, 0x9a, 0xdf, 0x8e, 0xc0, 0xa9, 0x89, 0xac, 0xc0, 0x16, 0x24, 0xcd, 0x5c, 0x0, 0xac, 0x86, 0x6d, 0xea}
// 		actual, err := x.PubKeyBytes()

// 		assert.NoError(t, err)
// 		assert.Equal(t, expected, actual)
// 	})
// }

func TestPayload_Verify(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		p := buildExamplePayload()
		assert.NoError(t, p.Verify())
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
