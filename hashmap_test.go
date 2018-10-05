package hashmap_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/nomasters/hashmap"

	"github.com/multiformats/go-multihash"
	"golang.org/x/crypto/nacl/sign"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	badData = "8J+ZifCfmYjwn5mK" // base64 encoded monkeys
)

func TestNewPayloadFromReader(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		actual, err := hashmap.NewPayloadFromReader(strings.NewReader(exampleValidPayload))
		assert.NoError(t, err)
		assert.Equal(t, examplePayload, actual)
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()

		shouldBeNil, err := hashmap.NewPayloadFromReader(&errorReader{})
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with invalid JSON", func(t *testing.T) {
		t.Parallel()

		exampleReader := strings.NewReader("not json lol")
		shouldBeNil, err := hashmap.NewPayloadFromReader(exampleReader)
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

		shouldBeNil, err := hashmap.NewPayloadFromReader(exampleReader)
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestNewValidator(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		v, err := examplePayload.NewValidator()

		assert.NoError(t, v.Validate())
		assert.NoError(t, err)
	})

	T.Run("with pubkey failure", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
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

		example := &hashmap.Payload{
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

		example := &hashmap.Payload{
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

		example := &hashmap.Payload{
			Data:      badData,
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with invalid signature method", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
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

		_, privKey, _ := sign.GenerateKey(rand.Reader)

		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: base64.StdEncoding.EncodeToString(privKey[31:]),
		}

		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestPayload_PubKeyBytes(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
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
		p := &hashmap.Payload{}
		json.Unmarshal([]byte(badPayload), &p)

		assert.Error(t, p.Verify())
	})
}

func TestPayload_GetData(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		actual, err := examplePayload.DataBytes()
		assert.NotEmpty(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error decoding data", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		actual, err := example.GetData()
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestData_ValidateTTL(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().UnixNano()),
			TTL:       hashmap.DataTTLMax,
		}
		actual := example.ValidateTTL()
		assert.NoError(t, actual)
	})

	T.Run("sets a valid TTL when not provided with one", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().UnixNano()),
		}
		actual := example.ValidateTTL()
		assert.NoError(t, actual)
	})

	T.Run("with too long a TTL", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{TTL: 1<<63 - 1}
		assert.Error(t, example.ValidateTTL())
	})

	T.Run("with an exceeded TTL", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().Nanosecond()),
			TTL:       int64(1 * time.Nanosecond),
		}
		actual := example.ValidateTTL()

		assert.Error(t, actual)
	})
}

func TestData_ValidateMessageSize(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		d, err := examplePayload.GetData()
		assert.NoError(t, err)

		assert.NoError(t, d.ValidateMessageSize())
	})

	T.Run("with error getting message bytes", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoibm90IGEgcmVhbCBiYXNlNjQgdGhpbmcgbG9sIiwidGltZXN0YW1wIjoxNTM0NDc3MjMwLCJ0dGwiOjg2NDAwLCJzaWdNZXRob2QiOiJuYWNsLXNpZ24tZWQyNTUxOSIsInZlcnNpb24iOiIwLjAuMSJ9",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		d, err := example.GetData()
		assert.NoError(t, err)

		assert.Error(t, d.ValidateMessageSize())
	})

	T.Run("with too much data", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      tooMuchData,
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		d, err := example.GetData()
		assert.NoError(t, err)

		assert.Error(t, d.ValidateMessageSize())
	})
}

func TestData_ValidateTimeStamp(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)
		assert.NoError(t, d.ValidateTimeStamp())
	})

	T.Run("with negative time difference", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)

		d.Timestamp = time.Now().Add(1 * time.Hour).UnixNano()
		assert.Error(t, d.ValidateTimeStamp())
	})

	T.Run("with too much submission drift", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)

		d.Timestamp = time.Now().Add(-24 * time.Hour).UnixNano()
		assert.Error(t, d.ValidateTimeStamp())
	})
}

func TestMultiHashToString(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		expected := "2Drjgb9YQZNX4C2X5iADiSprs4N3LCZBTy6GcnWQ83aFHoKjwg"
		actual := hashmap.MultiHashToString([]byte("example"))
		assert.Equal(t, expected, actual)
	})
}

func TestValidateMultiHash(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		example := hashmap.MultiHashToString([]byte("example"))
		assert.NoError(t, hashmap.ValidateMultiHash(example))
	})

	T.Run("with error validating hash string", func(t *testing.T) {
		t.Parallel()

		assert.Error(t, hashmap.ValidateMultiHash("not a hash, lol"))
	})

	T.Run("wrong code", func(t *testing.T) {
		t.Parallel()

		mh, err := multihash.Sum([]byte("here are thirty two characters!!"), 0, -1)
		require.NoError(t, err)

		assert.Error(t, hashmap.ValidateMultiHash(mh.B58String()))
	})
}
