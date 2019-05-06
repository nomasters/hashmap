package sig

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/nacl/sign"
)

// NaClSign holds a pointer to a 64 byte array used by NaCl Sign. It implements the
// Signer interface.
type NaClSign struct {
	PrivateKey *[64]byte
}

// GenNaclSign returns a randomly generated ed25519 private key in a NaClSign
func GenNaclSign() NaClSign {
	_, k, _ := sign.GenerateKey(rand.Reader)
	return NewNaClSign(k[:])
}

// NewNaClSign takes a ed25519 private key as a byte array amd returns a NaClSign
func NewNaClSign(privateKey []byte) NaClSign {
	var pk [64]byte
	copy(pk[:], privateKey)
	return NaClSign{
		PrivateKey: &pk,
	}
}

// Sign takes a message and returns a Bundle signed with a private key using NaCl Sign.
func (s NaClSign) Sign(message []byte) (Bundle, error) {
	var pubkey [32]byte
	copy(pubkey[:], s.PrivateKey[32:])

	sig := sign.Sign(nil, message, s.PrivateKey)[:sign.Overhead]

	b := Bundle{
		Alg: AlgNaClSign,
		Pub: pubkey[:],
		Sig: sig,
	}

	if ok := VerifyNaclSign(message, b); !ok {
		return Bundle{}, errors.New("verification sanity check failed on sign")
	}

	return b, nil
}

// VerifyNaclSign takes a message and a Bundle and bool indicating if the message
// is verified ny the signature.
func VerifyNaclSign(message []byte, b Bundle) bool {
	signedMessage := append(b.Sig, message...)
	var pubkey [32]byte
	copy(pubkey[:], b.Pub)
	_, ok := sign.Open(nil, signedMessage, &pubkey)
	return ok
}
