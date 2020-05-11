package sig

import (
	"errors"
	"sync"

	xmss "github.com/danielhavir/go-xmss"
)

// XMSS10 holds a pointer to a 132 byte array used by XMSS. It implements the
// Signer interface.
type XMSS10 struct {
	m          sync.RWMutex
	PrivateKey [132]byte
}

// GenXMSS10 returns a randomly generated XMSS SHA2_10_256 private key in a NaClSign
func GenXMSS10() *XMSS10 {
	k, _ := xmss.GenerateXMSSKeypair(xmss.SHA2_10_256)
	return NewXMSS10((*k)[:])
}

// NewXMSS10 takes an XMSS private key as a byte array and returns a NaClSign
func NewXMSS10(privateKey []byte) *XMSS10 {
	var pk [132]byte
	copy(pk[:], privateKey)
	return &XMSS10{
		PrivateKey: pk,
	}
}

// Sign takes a message and returns a Bundle signed with a private key using XMSS SHA2_10_256.
func (s *XMSS10) Sign(message []byte) (Bundle, error) {
	s.m.Lock()
	prv := xmss.PrivateXMSS(s.PrivateKey[:])
	pub := make([]byte, 64)
	copy(pub[:32], prv[100:])
	copy(pub[32:], prv[68:100])

	prm := xmss.SHA2_10_256
	sig := *prv.Sign(prm, message)
	copy(s.PrivateKey[:], prv)
	s.m.Unlock()

	b := Bundle{
		Alg: AlgXMSS10,
		Pub: pub,
		Sig: Bytes(sig[:prm.SignBytes()]),
	}

	if ok := VerifyXMSS10(message, b); !ok {
		return Bundle{}, errors.New("verification sanity check failed on sign")
	}

	return b, nil
}

// VerifyXMSS10 takes a message and a Bundle and bool indicating if the message
// is verified ny the signature.
func VerifyXMSS10(msg []byte, b Bundle) bool {
	prm := xmss.SHA2_10_256
	sig := append(b.Sig, msg...)
	m := make([]byte, prm.SignBytes()+len(msg))
	return xmss.Verify(prm, m, sig, []byte(b.Pub))
}
