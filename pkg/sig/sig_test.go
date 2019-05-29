package sig

import (
	"testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	m := []byte("sign me, plz.")

	t.Run("naclSign", func(t *testing.T) {
		s := GenNaclSign()
		b, _ := s.Sign(m)
		if ok := Verify(m, b); !ok {
			t.Error("verification failed")
		}
	})
	t.Run("invalidAlg", func(t *testing.T) {
		b := Bundle{
			Alg: 0,
		}
		if ok := Verify(m, b); ok {
			t.Error("verification should fail")
		}
	})
}
