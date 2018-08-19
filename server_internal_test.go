package hashmap

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	badData = "8J+ZifCfmYjwn5mK" // base64 encoded monkeys
)

func buildExamplePayload(t *testing.T, msg string) *Payload {
	t.Helper()
	pk := GenerateKey()
	opts := Options{}
	p, err := GeneratePayload(opts, pk)
	require.NoError(t, err)
	return p
}

var _ Storage = (*badStorage)(nil)

type badStorage struct{}

func (b *badStorage) Get(string) (PayloadWithMetadata, error) {
	return PayloadWithMetadata{}, errors.New("arbitrary")
}
func (b *badStorage) Set(string, PayloadWithMetadata) error { return errors.New("arbitrary") }
func (b *badStorage) Delete(string) error                   { return errors.New("arbitrary") }

func TestBuildSubmitHandler(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		s, err := NewStorage(MemoryStorage, nil)
		require.NoError(t, err)

		exampleBody, err := json.Marshal(buildExamplePayload(t, "whatever"))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(exampleBody))
		res := httptest.NewRecorder()

		buildSubmitHandler(s)(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.NotEmpty(t, res.Body.String())
	})

	T.Run("with too large a message", func(t *testing.T) {
		t.Skip() // TODO: figure this out

		s, err := NewStorage(MemoryStorage, nil)
		require.NoError(t, err)

		p := buildExamplePayload(t, "")

		d := Data{
			Message:   tooMuchData,
			Timestamp: time.Now().Unix(),
			TTL:       DataTTLMax,
			SigMethod: DefaultSigMethod,
			Version:   Version,
		}
		data, err := json.Marshal(d)
		p.Data = base64.StdEncoding.EncodeToString(data)

		body, err := json.Marshal(&p)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(body))
		res := httptest.NewRecorder()

		buildSubmitHandler(s)(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})

	T.Run("with an invalid TTL", func(t *testing.T) {
		t.Parallel()

		s, err := NewStorage(MemoryStorage, nil)
		require.NoError(t, err)

		exampleBody, err := json.Marshal(examplePayload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(exampleBody))
		res := httptest.NewRecorder()

		buildSubmitHandler(s)(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, res.Body.String(), "ttl exceeded\n")
	})

	T.Run("with invalid body", func(t *testing.T) {
		t.Parallel()

		s, err := NewStorage(MemoryStorage, nil)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", strings.NewReader(`lol this is bad`))
		res := httptest.NewRecorder()

		buildSubmitHandler(s)(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, res.Body.String(), "invalid payload\n")
	})

	T.Run("with invalid storage", func(t *testing.T) {
		t.Parallel()

		bs := &badStorage{}

		exampleBody, err := json.Marshal(buildExamplePayload(t, "whatever"))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(exampleBody))
		res := httptest.NewRecorder()

		buildSubmitHandler(bs)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, res.Body.String(), "internal error saving payload\n")
	})
}

func TestGetPayloadHandleFunc(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleBody, err := json.Marshal(buildExamplePayload(t, "whatever"))
		require.NoError(t, err)

		testReq := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(exampleBody))
		res := httptest.NewRecorder()
		req := testReq.WithContext(context.WithValue(context.Background(), payloadCtxKey, Payload{}))

		getPayloadHandleFunc(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, `{"data":"","sig":"","pubkey":""}`, res.Body.String())
	})
}

func testHandler(called *bool) func(http.Handler) http.Handler {
	return func(http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, req *http.Request) {
			*called = true
		}
		return http.HandlerFunc(fn)
	}
}

func TestBuildPrivateKeyHashMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		examplePayload := buildExamplePayload(t, "hey")
		pk, err := examplePayload.PubKeyBytes()
		require.NoError(t, err)

		hash := MultiHashToString(pk)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost/%s", hash), nil)
		res := httptest.NewRecorder()

		s := NewMemoryStore()
		s.internal = map[string]PayloadWithMetadata{hash: {Payload: *examplePayload}}
		r := buildRouter(s)
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		//assert.Equal(t, ``, res.Body.String())
	})
}
