package sig

import (
	"bytes"
	"testing"
)

func TestSign(t *testing.T) {
	t.Parallel()

	m := []byte("sign me, plz.")
	s := GenNaclSign()

	t.Run("normal", func(t *testing.T) {
		if _, err := s.Sign(m); err != nil {
			t.Error(err)
		}
	})
	t.Run("malformed", func(t *testing.T) {
		badKey := bytes.Repeat([]byte{5}, 64)
		copy(s.PrivateKey[:], badKey)
		if _, err := s.Sign(m); err == nil {
			t.Error("malformed signature not caught")
		}
	})
}
