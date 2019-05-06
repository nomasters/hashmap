package payload

import (
	"testing"
	"time"
)

func TestNewPayload(t *testing.T) {
	t.Parallel()
	p := Payload{
		Version:   V1,
		Timestamp: time.Now(),
		TTL:       30 * time.Second,
		SigSets: []SigSet{
			SigSet{
				Alg: NACLSign,
				Pub: []byte("so fake"),
				Sig: []byte("yep, the fakest"),
			},
		},
		Data: []byte("hello, world"),
	}

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

}
