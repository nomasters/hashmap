package payload

import (
	"encoding/base64"
	"testing"
	"time"

	sig "github.com/nomasters/hashmap/pkg/sig"
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

func TestPayloadMethods(t *testing.T) {
	t.Parallel()
	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())
	now := time.Now()
	message := []byte("hello, world")
	t.Run("validPayloadSize Bad Marshal", func(t *testing.T) {
		bad := time.Unix(-99999999999, 0)
		p, err := Generate(message, signers, WithTimestamp(bad))
		if err != nil {
			t.Error(err)
		}
		if p.validPayloadSize() {
			t.Error("failed to catch marshalling error")
		}
	})
	t.Run("withinSubmitWindow negative diff", func(t *testing.T) {
		p, err := Generate(message, signers, WithTimestamp(now))
		if err != nil {
			t.Error(err)
		}
		if !p.withinSubmitWindow(now.Add(-3 * time.Second)) {
			t.Fail()
		}
		if p.withinSubmitWindow(now.Add(-10 * time.Second)) {
			t.Fail()
		}
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()
	var signers []sig.Signer
	s1 := sig.GenNaclSign()
	signers = append(signers, s1)
	now := time.Now()
	message := []byte("hello, world")
	p, _ := Generate(message, signers, WithTimestamp(now))

	t.Run("Normal Operation", func(t *testing.T) {
		if err := validate(p); err != nil {
			t.Error(err)
		}
	})
	t.Run("Endpoint", func(t *testing.T) {
		p, _ := Generate(message, signers, WithTimestamp(now))
		e := base64.URLEncoding.EncodeToString(p.PubKeyHash())
		if err := validate(p, WithValidateEndpoint(e)); err != nil {
			t.Error(err)
		}
		if err := validate(p, WithValidateEndpoint("BAD_ENDPOINT")); err == nil {
			t.Error("validate did not catch MaxMessageSize")
		}
		if err := validate(p, WithValidateEndpoint("")); err != nil {
			t.Error(err)
		}
	})
	t.Run("Data Size", func(t *testing.T) {
		message := make([]byte, MaxMessageSize+1)
		p, _ := Generate(message, signers, WithTimestamp(now))
		if err := validate(p, WithValidateDataSize(true)); err == nil {
			t.Error("validate did not catch MaxMessageSize")
		}
		if err := validate(p, WithValidateDataSize(false)); err != nil {
			t.Error(err)
		}
	})
	t.Run("Payload Size", func(t *testing.T) {
		message := make([]byte, MaxPayloadSize+1)
		p, _ := Generate(message, signers, WithTimestamp(now))
		if err := validate(p,
			WithValidatePayloadSize(true),
			WithValidateDataSize(false)); err == nil {
			t.Error("validate did not catch MaxPayloadSize")
		}
		if err := validate(p,
			WithValidatePayloadSize(false),
			WithValidateDataSize(false)); err != nil {
			t.Error(err)
		}
	})
	t.Run("Version", func(t *testing.T) {
		p, _ := Generate(message, signers,
			WithTimestamp(now),
			WithVersion(Version(3)))
		if err := validate(p, WithValidateVersion(true)); err == nil {
			t.Error("validate did not catch bad version")
		}
		if err := validate(p, WithValidateVersion(false)); err != nil {
			t.Error(err)
		}
	})
	t.Run("Future", func(t *testing.T) {
		p, _ := Generate(message, signers, WithTimestamp(now.Add(MaxSubmitWindow+1)))
		if err := validate(p,
			WithValidateFuture(true),
			WithReferenceTime(now)); err == nil {
			t.Error("validate did not catch timestamp in future")
		}
		if err := validate(p,
			WithValidateFuture(false),
			WithReferenceTime(now)); err != nil {
			t.Error(err)
		}
	})
	t.Run("TTL", func(t *testing.T) {
		p, _ := Generate(message, signers, WithTTL(MaxTTL+1))
		if err := validate(p, WithValidateTTL(true)); err == nil {
			t.Error("validate did not catch bad ttl")
		}
		if err := validate(p, WithValidateTTL(false)); err != nil {
			t.Error(err)
		}
	})
	t.Run("expiration", func(t *testing.T) {
		p, _ := Generate(
			message,
			signers,
			WithTTL(10*time.Second),
			WithTimestamp(now.Add(-11*time.Second)))
		if err := validate(p, WithValidateExpiration(true)); err == nil {
			t.Error("validate did not catch expired ttl")
		}
		if err := validate(p, WithValidateExpiration(false)); err != nil {
			t.Error(err)
		}
	})
}
