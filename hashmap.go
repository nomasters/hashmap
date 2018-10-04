package hashmap

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/multiformats/go-multihash"
)

// A set of constants used to set defaults and Max settings
const (
	DataTTLDefault   = 86400  // 1 day in seconds
	DataTTLMax       = 604800 // 1 week in seconds
	DefaultSigMethod = "nacl-sign-ed25519"
	Version          = "0.0.1"
	MaxPostBodySize  = 2000 // 2KB
	MaxSubmitDrift   = 15 * time.Second
	MaxMessageBytes  = 512
	Blake2b256Code   = 45600
)

// Payload is the primary wrapper struct for HashMap Values and submissions
type Payload struct {
	Data      string `json:"data"`
	Signature string `json:"sig"`
	PublicKey string `json:"pubkey"`
}

// PayloadWithMetadata is a struct that contains a Payload and related Metadata.
// The Metadata is is a map[string]interface{} to give flexibility to systems
// that use it.
type PayloadWithMetadata struct {
	Payload
	Metadata map[string]string
}

// Data is the struct for the Data in a Payload. It contains all data that is
// signed by the Payload Pubkey
type Data struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	TTL       int64  `json:"ttl"`
	SigMethod string `json:"sigMethod"`
	Version   string `json:"version"`
}

// NewPayloadFromReader returns a fully verified Payload from an io.Reader source.
// This includes verifying signature and size restrictions.
func NewPayloadFromReader(r io.Reader) (*Payload, error) {
	p := Payload{}
	// read the body with strict limit on body size
	limitedReader := &io.LimitedReader{R: r, N: MaxPostBodySize}
	body, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		log.Printf("error reading payload: %v\n", err)
		return nil, errors.New("error reading payload")
	}

	// unmarshal the payload, returning an error if it fails
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("invalid payload: %v\n", err)
		return nil, errors.New("invalid payload")
	}
	if err := p.Verify(); err != nil {
		return nil, err
	}

	return &p, nil
}

// Verify validates a message signature and enforces message requirements
func (p Payload) Verify() error {
	v, err := p.NewValidator()
	if err != nil {
		return err
	}
	return v.Validate()
}

// NewValidator decodes the pubKey, Signature, and DataBytes to byte slices.
// It then Unmarshals the Data and analyzes the SigMethod, either returning
// a Validator or an error.
func (p Payload) NewValidator() (Validator, error) {
	pubKey, err := p.PubKeyBytes()
	if err != nil {
		return nil, err
	}
	sig, err := p.SignatureBytes()
	if err != nil {
		return nil, err
	}
	dataBytes, err := p.DataBytes()
	if err != nil {
		return nil, err
	}
	data, err := p.GetData()
	if err != nil {
		return nil, err
	}

	switch data.SigMethod {
	case "nacl-sign-ed25519":
		if len(pubKey) != 32 {
			return nil, errors.New("invalid pubKey length")
		}
		return NewNaClSignEd25519(append(sig, dataBytes...), pubKey), nil
	default:
		return nil, errors.New("invalid signature method")
	}

}

// PubKeyBytes method decodes a Payload.PublicKey and returns a slice of bytes and an error
func (p Payload) PubKeyBytes() ([]byte, error) {
	pubKey, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		log.Printf("invalid pubKey encoding: %v\n", err)
		return pubKey, errors.New("invalid pubKey encoding, expecting base64")
	}
	return pubKey, nil
}

// SignatureBytes method decodes a Payload.Signature and returns a slice of bytes and an error
func (p Payload) SignatureBytes() ([]byte, error) {
	sig, err := base64.StdEncoding.DecodeString(p.Signature)
	if err != nil {
		log.Printf("invalid signature encoding: %v\n", err)
		return sig, errors.New("invalid signature encoding, expecting base64")
	}
	return sig, nil
}

// DataBytes method decodes a Payload.Data and returns a slice of bytes and an error
func (p Payload) DataBytes() ([]byte, error) {
	d, err := base64.StdEncoding.DecodeString(p.Data)
	if err != nil {
		log.Printf("invalid data encoding: %v\n", err)
		return nil, errors.New("invalid data encoding, expecting base64")
	}
	return d, nil
}

// GetData method decodes and unmarshals a Payload.Data and returns a pointer to data and an error
func (p Payload) GetData() (*Data, error) {
	// decode data
	data, err := p.DataBytes()
	if err != nil {
		return nil, err
	}

	d := Data{}
	if err := json.Unmarshal(data, &d); err != nil {
		log.Printf("invalid data: %v\n", err)
		return nil, errors.New("invalid data")
	}

	return &d, nil
}

// ValidateTTL checks that a TTL is configured within the boundries of a proper TTL
// and then checks the TTL against the diff of the timestamp & time.Now().
// if any of the checks fail, ValidateTTL returns an error
func (d Data) ValidateTTL() error {
	t := d.TTL

	if t > DataTTLMax {
		return fmt.Errorf("message ttl exceeds max allowed of %v\n", DataTTLMax)
	}

	if t == 0 {
		t = DataTTLDefault
	}

	// convert to duration
	ttl := time.Duration(t) * time.Second

	timeStamp := time.Unix(secondsAndNanoseconds(d.Timestamp))
	now := time.Now()
	diff := now.Sub(timeStamp)

	if diff > ttl {
		return errors.New("ttl exceeded")
	}

	return nil
}

// ValidateMessageSize decodes data and checks that it does not excede MaxMessageBytes.
// ValidateMessageSize returns an error if any validation fails.
func (d Data) ValidateMessageSize() error {
	data, err := d.MessageBytes()
	if err != nil {
		return err
	}

	if len(data) > MaxMessageBytes {
		return fmt.Errorf("message exceeds max allowed: %v\n", MaxMessageBytes)
	}
	return nil
}

// ValidateTimeStamp compares time.Now to message Timestamp. If the difference
// exceeds MaxSubmitDrift, it returns an error. This is to prevent replay attacks.
func (d Data) ValidateTimeStamp() error {
	timeStamp := time.Unix(secondsAndNanoseconds(d.Timestamp))
	now := time.Now()
	diff := now.Sub(timeStamp)

	// get absolute value of time difference
	if diff.Seconds() < 0 {
		diff = -diff
	}

	if diff > MaxSubmitDrift {
		return errors.New("max submission time drift exceeded for message")
	}

	return nil
}

// MessageBytes method decodes a Message.Data and returns a slice of bytes and an error
func (d Data) MessageBytes() ([]byte, error) {
	m, err := base64.StdEncoding.DecodeString(d.Message)
	if err != nil {
		log.Printf("invalid message encoding: %v\n", err)
		return m, errors.New("invalid message encoding, expecting base64")
	}
	return m, nil
}

// MultiHashToString takes a slice of bytes, shahes to blake2b-256
// and returns a BTC/IPFS style Base58 encoded string
func MultiHashToString(b []byte) string {
	// TODO: In the future, this should be version number aware
	mh, _ := multihash.Sum(b, Blake2b256Code, -1)
	return mh.B58String()
}

// ValidateMultiHash takes a multihash encoded in base58, decodes, and validates
// against the valid results. This may change over time, if we support more hashes.
func ValidateMultiHash(hash string) error {
	mh, err := multihash.FromB58String(hash)
	if err != nil {
		log.Printf("%v failed to decode multihash with error: %s\n", hash, err)
		return errors.New("multiHash Decode failed")
	}

	// the `multihash.FromB58String` call above calls this function and returns the
	//  error if it's not nil, so we can safely ignore the error here
	dh, _ := multihash.Decode(mh)

	if dh.Length != 32 {
		return errors.New("multiHash length invalid")
	}
	if int(dh.Code) != Blake2b256Code {
		return errors.New("multiHash code invalid")
	}
	if len(dh.Digest) != 32 {
		return errors.New("pubKey hash length invalid")
	}
	return nil
}

// secondsAndNanoseconds takes a full time UnixNano int64 and returns the seconds and the nanoseconds
// expected by the Unix() parser for the time library.
func secondsAndNanoseconds(i int64) (s, n int64) {
	s = i / 1000000000
	n = i - (s * 1000000000)
	return
}
