package payload

import (
	"testing"
	"time"

	sig "github.com/nomasters/hashmap/sig"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())
	message := []byte("hello, world")
	now := time.Now()
	p, _ := Generate(message, signers, WithTimestamp(now))

	t.Run("Normal Operation", func(t *testing.T) {
		if err := p.Verify(); err != nil {
			t.Error(err)
		}
	})
	t.Run("Failed Validation", func(t *testing.T) {
		err := p.Verify(
			WithReferenceTime(now.Add(15*time.Second)),
			WithServerMode(true),
		)
		if err == nil {
			t.Error("Verify did not catch invalid payload")
		}
	})
	t.Run("Failed Verification", func(t *testing.T) {
		p.SigBundles[0].Sig = []byte("bad_bytes")
		err := p.Verify()
		if err == nil {
			t.Error("Verify did not catch invalid signature")
		}
	})
}
