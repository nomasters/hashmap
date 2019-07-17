package sigutil

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	blake2b "golang.org/x/crypto/blake2b"
	sig "github.com/nomasters/hashmap/pkg/sig"
)

func init() {
	gob.Register(&sig.NaClSign{})
	gob.Register(&sig.XMSS10{})
}

// NewDefaultSigners returns a slice of sig.Signer that
// includes an ed25519 sig
func NewDefaultSigners() (s []sig.Signer) {
	return append(s, sig.GenNaclSign())
}

// NewExperimentalSigners returns a slice of sig.Signer that
// includes an ed25519 sig and an XMSS sig
func NewExperimentalSigners() (s []sig.Signer) {
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

// SignAll is a helper function for quickly signing a byteslice with a group
// of signers. It takes a slice of sig.Signers and the bytes for signing and
// returns a slice of sig.Bundles and an error
func SignAll(message []byte, signers []sig.Signer) ([]sig.Bundle, error){
	var sigBundles []sig.Bundle
	for _, s := range signers {
		bundle, err := s.Sign(message)
		if err != nil {
			return nil, err
		}
		sigBundles = append(sigBundles, bundle)
	}
	return sigBundles, nil
}

// VerifyAll takes message bytes and a slice of sig.Bundles and returns a boolean
// value of true if all signatures verify, otherwise it returns false
func VerifyAll(message []byte, bundles []sig.Bundle) bool {
	for _, bundle := range bundles {
		if !sig.Verify(message, bundle) {
			return false
		}
	}
	return true
}

// BundlePubKeys returns a byte slice of all pubkeys concatenated in the index
// order of the slice of sig.Bundles.
func BundlePubKeys(bundles []sig.Bundle) []byte {
	var o []byte
	for _, b := range bundles {
		o = append(o, b.Pub...)
	}
	return o
}

// BundleHash returns a byte slice of the blake2b-512 hash of PubKeys
func BundleHash(bundles []sig.Bundle) []byte {
	b := blake2b.Sum512(BundlePubKeys(bundles))
	return b[:]
}

// EncodedBundleHash returns a url-safe base64 encoded endpoint string of PubKeyHash
func EncodedBundleHash(bundles []sig.Bundle) string {
	return base64.URLEncoding.EncodeToString(BundleHash(bundles))
}