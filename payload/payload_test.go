package payload

import (
	"testing"

	"github.com/nomasters/hashmap/sig"
)

func TestNewPayload(t *testing.T) {
	t.Parallel()

	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())

	message := []byte("hello, world")

	p, err := Generate(message, signers)
	if err != nil {
		t.Error(err)
	}

	t.Log(p)

	encoded, err := Marshal(p)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%x", encoded)
	decoded, err := Unmarshal(encoded)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", string(decoded.Data))
	t.Logf("%v", decoded.Timestamp)
	t.Logf("%v", decoded.TTL)

}
