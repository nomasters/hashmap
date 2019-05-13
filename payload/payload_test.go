package payload

import (
	"testing"
	"time"

	"github.com/nomasters/hashmap/sig"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	var signers []sig.Signer
	m := []byte("hello, world")

	t.Run("Normal Operation", func(t *testing.T) {
		s := append(signers, sig.GenNaclSign())
		if _, err := Generate(m, s); err != nil {
			t.Error(err)
		}
	})
	t.Run("With Options", func(t *testing.T) {
		s := append(signers, sig.GenNaclSign())
		ttl := 5 * time.Second
		timestamp := time.Now().Add(5 * time.Minute)
		version := Version(2)

		p, err := Generate(m, s,
			WithTTL(ttl),
			WithTimestamp(timestamp),
			WithVersion(version),
		)
		if err != nil {
			t.Error(err)
		}
		if p.TTL != ttl {
			t.Error("ttl mismatch")
		}
		if p.Timestamp != timestamp {
			t.Error("timestamp mismatch")
		}
		if p.Version != version {
			t.Error("version mismatch")
		}

	})
	t.Run("Empty Signers", func(t *testing.T) {
		if _, err := Generate(m, signers); err == nil {
			t.Error("should reject on missing signers")
		}
	})
	t.Run("Malformed Signer", func(t *testing.T) {
		s := append(signers, sig.NaClSign{PrivateKey: &[64]byte{}})
		if _, err := Generate(m, s); err == nil {
			t.Error("should reject a malformed signer")
		}
	})

}

func TestMarshal(t *testing.T) {
	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())
	message := []byte("hello, world")

	t.Run("Normal Operation", func(t *testing.T) {
		p, err := Generate(message, signers)
		if err != nil {
			t.Error(err)
		}
		if _, err := Marshal(p); err != nil {
			t.Error(err)
		}
	})

	t.Run("Invalid Timestamp", func(t *testing.T) {
		bad := time.Unix(-99999999999, 0)
		p, err := Generate(message, signers, WithTimestamp(bad))
		if err != nil {
			t.Error(err)
		}
		if _, err := Marshal(p); err == nil {
			t.Error("failed to catch invalid timestamp")
		}
	})

}

func TestUnmarshal(t *testing.T) {
	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())
	message := []byte("hello, world")
	p, err := Generate(message, signers)
	if err != nil {
		t.Error(err)
	}
	encoded, err := Marshal(p)
	if err != nil {
		t.Error(err)
	}
	if _, err := Unmarshal(encoded); err != nil {
		t.Error(err)
	}
}
