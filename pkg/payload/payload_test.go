package payload

import (
	"bytes"
	"testing"
	"time"

	sig "github.com/nomasters/hashmap/pkg/sig"
	blake2b "golang.org/x/crypto/blake2b"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	var signers []sig.Signer
	m := []byte("hello, world")

	t.Run("Normal Operation", func(t *testing.T) {
		// single sig
		s1 := append(signers, sig.GenNaclSign())
		if _, err := Generate(m, s1); err != nil {
			t.Error(err)
		}
		// multisig
		s2 := append(signers, sig.GenNaclSign(), sig.GenNaclSign())
		if _, err := Generate(m, s2); err != nil {
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
		bad := &sig.NaClSign{PrivateKey: [64]byte{}}
		s := append(signers, bad)
		if _, err := Generate(m, s); err == nil {
			t.Error("should reject a malformed signer")
		}
	})

}

func TestPayloadCoreMethods(t *testing.T) {
	t.Parallel()

	// PubKeyBytes tests that the PubKeyBytes method properly generates a concat of
	// all public keys in proper order
	t.Run("PubKeyBytes", func(t *testing.T) {
		var signers []sig.Signer
		sig1, sig2, sig3 := sig.GenNaclSign(), sig.GenNaclSign(), sig.GenNaclSign()
		signers = append(signers, sig1, sig2, sig3)
		message := []byte("")
		p, err := Generate(message, signers)
		if err != nil {
			t.Error(err)
		}
		var naclSigns []*sig.NaClSign
		var refConcat []byte
		naclSigns = append(naclSigns, sig1, sig2, sig3)
		for _, s := range naclSigns {
			pubkeyBytes := s.PrivateKey[32:]
			refConcat = append(refConcat, pubkeyBytes...)
		}
		if !bytes.Equal(refConcat, p.PubKeyBytes()) {
			t.Error("concatenated pubkey bytes are invalid")
		}
	})
	t.Run("PubKeyHash", func(t *testing.T) {
		var signers []sig.Signer
		sig1 := sig.GenNaclSign()
		signers = append(signers, sig1)
		message := []byte("")
		p, err := Generate(message, signers)
		if err != nil {
			t.Error(err)
		}
		b := blake2b.Sum512(sig1.PrivateKey[32:])
		if !bytes.Equal(b[:], p.PubKeyHash()) {
			t.Error("pubkey hash bytes are invalid")
		}
	})
}
