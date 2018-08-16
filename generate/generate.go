package generate

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/nomasters/hashmap"
	"golang.org/x/crypto/nacl/sign"
)

// Options is a struct used for passing in Message, TTL, and Timestamp options
// and is used in conjunction with CLI flags
type Options struct {
	Message   string
	TTL       int64
	Timestamp int64
}

// Payload takes an Options struct and the bytes of a private key
// and returns a json encoded payload and an error
func Payload(opts Options, pk []byte) ([]byte, error) {

	if opts.Message == "" {
		opts.Message = `{"content":"hello, world. This is data stored in HashMap."}`
	}

	if opts.TTL == 0 {
		opts.TTL = hashmap.DataTTLDefault
	}

	if opts.Timestamp == 0 {
		opts.Timestamp = time.Now().Unix()
	}

	message := base64.StdEncoding.EncodeToString([]byte(opts.Message))
	d := hashmap.Data{
		Message:   message,
		Timestamp: opts.Timestamp,
		TTL:       opts.TTL,
		SigMethod: hashmap.DefaultSigMethod,
		Version:   hashmap.Version,
	}
	data, err := json.Marshal(d)
	if err != nil {
		return []byte(""), err
	}

	var privateKey [64]byte
	var publicKey [32]byte
	copy(privateKey[:], pk)
	copy(publicKey[:], pk[32:])

	s := sign.Sign(nil, data, &privateKey)[:64]
	sig := base64.StdEncoding.EncodeToString(s)

	p := hashmap.Payload{
		Data:      base64.StdEncoding.EncodeToString(data),
		Signature: sig,
		PublicKey: base64.StdEncoding.EncodeToString(publicKey[:]),
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return []byte(""), err
	}

	if _, valid := sign.Open(nil, append(s, data...), &publicKey); !valid {
		return []byte(""), errors.New("signature failed to validate")
	}

	return payload, nil

}

// Key returns a randomly generated ed25519 private key in bytes
func Key() []byte {
	_, privKey, _ := sign.GenerateKey(rand.Reader)
	return privKey[:]
}
