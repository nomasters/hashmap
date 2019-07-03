package sig

import (
	"bytes"
	"testing"
)

func TestXMSS10(t *testing.T) {
	t.Parallel()
	m := []byte("sign me, plz.")
	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		s := GenXMSS10()
		pre := s.PrivateKey[3]
		if _, err := s.Sign(m); err != nil {
			t.Error(err)
		}
		post := s.PrivateKey[3]
		if pre >= post {
			t.Error("counter state invalid")
		}
	})
	t.Run("malformed", func(t *testing.T) {
		t.Parallel()
		s := GenXMSS10()
		badKey := bytes.Repeat([]byte{5}, 100)
		copy(s.PrivateKey[15:], badKey)
		if _, err := s.Sign(m); err == nil {
			t.Error("malformed signature not caught")
		}
	})
}
