package analyze

import (
	// sig "github.com/nomasters/hashmap/pkg/sig"
	sigutil "github.com/nomasters/hashmap/pkg/sig/sigutil"
)

type KeySet struct {
	Hash    string `json:"pubkey_hash"`
	Signers []Signer `json:"signers"`
	Valid   bool `json:"valid"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

type Signer struct {
	Type  string `json:"type"`
	Count string `json:"count,omitempty"`
	PQR   bool   `json:"pqr"`
}

func NewKeySet(b []byte) (*KeySet, error) {
	signers, err := sigutil.Decode(b)
	if err != nil {
		return nil, err
	}
	m := []byte("hello, world")
	sigBundles,  err := sigutil.SignAll(m, signers)
	if err != nil {
		return nil, err
	}

	// TODO:
	// write signer enum -> string function
	// write a function to check state on XMSS, should output XX of XXX
	// write a function for PQR metadata


	// get pubkey hash
	k := KeySet{
		Hash: sigutil.EncodedBundleHash(sigBundles),
		Valid: sigutil.VerifyAll(m, sigBundles),
	}
	return &k, nil
}