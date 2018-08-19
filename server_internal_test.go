package hashmap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildExamplePayload(t *testing.T, msg string) *Payload {
	t.Helper()
	pk := GenerateKey()
	opts := Options{}
	p, err := GeneratePayload(opts, pk)
	require.NoError(t, err)
	return p
}

func TestSubmitHandleFunc(T *testing.T) {
	T.Skip()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		exampleBody, err := json.Marshal(buildExamplePayload(t, "whatever"))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(exampleBody))
		res := httptest.NewRecorder()

		submitHandleFunc(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func TestGetPayloadHandleFunc(T *testing.T) {
	T.Skip()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

	})
}

func TestPKHashCtx(T *testing.T) {
	T.Skip()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

	})
}
