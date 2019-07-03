package sigutil

import (
	"bytes"
	"encoding/gob"

	sig "github.com/nomasters/hashmap/pkg/sig"
)

func init() {
	gob.Register(&sig.NaClSign{})
	gob.Register(&sig.XMSS10{})
}

// NewDefaultSigners returns a slice of sig.Signer that
// includes an ed25519 sig and xmss sha2_10_256 sig
func NewDefaultSigners() (s []sig.Signer) {
	return append(s, sig.GenNaclSign(), sig.GenXMSS10())
}

// Encode takes a slice of sig.Signer and returns a gob
// encoded byte slice and an error
func Encode(s []sig.Signer) ([]byte, error) {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(s)
	return buffer.Bytes(), err
}

// Decode takes a gob encoded byte slice and returns
// a []sig.Signer and an error.
func Decode(b []byte) (s []sig.Signer, err error) {
	var buffer bytes.Buffer
	buffer.Write(b)
	err = gob.NewDecoder(&buffer).Decode(&s)
	return s, err
}