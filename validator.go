package hashmap

import (
	"errors"
	"log"

	"golang.org/x/crypto/nacl/sign"
)

// Validator is an interface that takes a Validate method.
// each supported `sigMethod` will have a corresponding
// Validator for the payload
type Validator interface {
	Validate() error
}

// NaClSignEd25519 is the struct used by the Validator for the sigMethod: `nacl-sign-ed25519`
type NaClSignEd25519 struct {
	SignedMessage []byte
	PublicKey     *[32]byte
}

// NewNaClSignEd25519 takes to byte slices and returns a pointer to NaClSignEd25519 struct
func NewNaClSignEd25519(s, p []byte) *NaClSignEd25519 {
	var pk [32]byte
	copy(pk[:], p)
	return &NaClSignEd25519{
		SignedMessage: s,
		PublicKey:     &pk,
	}
}

// Validate conforms to the Validator interface and checks the validity of the NaClSignEd25519
// signed mesage against the Ed25519 pubkey
func (n NaClSignEd25519) Validate() error {
	// verify signature
	_, valid := sign.Open(nil, append(n.SignedMessage), n.PublicKey)
	if !valid {
		log.Printf("invalid signature: %x\n", n.PublicKey)
		return errors.New("invalid signature")
	}
	return nil
}
