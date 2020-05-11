package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	sig "github.com/nomasters/hashmap/pkg/sig"
	sigutil "github.com/nomasters/hashmap/pkg/sig/sigutil"
)

// Version type is used for setting the hashmap implementation version.
type Version uint16

// Bytes is a byte slice with special json encoding properties
type Bytes []byte

const (
	// V0 is deprecated and should be deemed invalid
	V0 Version = iota
	// V1 is the current version of the payload spec
	V1
)

const (
	// DefaultTTL is set to 24 hours
	DefaultTTL = 24 * time.Hour
)

var (
	defaultVersion = V1
)

// Payload holds all information related to a Hashmap Payload that will be handled
// for signing and validation. This struct is used by both client and server and
// includes all necessary methods for encoding, decoding, signing, an verifying itself.
type Payload struct {
	Version    Version       `json:"version"`
	Timestamp  time.Time     `json:"timestamp"`
	TTL        time.Duration `json:"ttl"`
	SigBundles []sig.Bundle  `json:"sig_bundles"`
	Data       Bytes         `json:"data"`
}

// Option is used for interacting with Context when setting options for Generate and Verify
type Option func(*options)

// options contains private fields used for Option
type options struct {
	version   Version
	timestamp time.Time
	ttl       time.Duration
	validate  validateContext
}

// WithVersion takes a Version and returns an Option
func WithVersion(v Version) Option {
	return func(o *options) {
		o.version = v
	}
}

// WithTimestamp takes a time.Time and returns an Option
func WithTimestamp(t time.Time) Option {
	return func(o *options) {
		o.timestamp = t
	}
}

// WithTTL takes a time.Duration and returns an Option
func WithTTL(d time.Duration) Option {
	return func(o *options) {
		o.ttl = d
	}
}

// parseOptions takes a arbitrary number of Option funcs and returns a Context with defaults
// for version, timestamp, and ttl, and validate rules.
func parseOptions(opts ...Option) options {
	now := time.Now()

	o := options{
		version:   defaultVersion,
		timestamp: now,
		ttl:       DefaultTTL,
		validate: validateContext{
			ttl:           true,
			expiration:    true,
			payloadSize:   true,
			dataSize:      true,
			version:       true,
			submitTime:    false,
			futureTime:    true,
			referenceTime: now,
		},
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// Generate takes a message, signers, and a set of options and returns a payload or error.
// This function defaults to time.Now() and the default TTL of 24 hours. Generate Requires
// at least one signer, but can sign with many signers. Sort order is important though, The unique
// order of the signers pubkeys are what is responsible for generating the endpoint hash.
func Generate(message []byte, signers []sig.Signer, opts ...Option) (Payload, error) {
	if len(signers) == 0 {
		return Payload{}, errors.New("Generate must have at least one signer")
	}
	o := parseOptions(opts...)
	p := Payload{
		Version:   o.version,
		Timestamp: o.timestamp,
		TTL:       o.ttl,
		Data:      message,
	}

	sigBundles, err := sigutil.SignAll(p.SigningBytes(), signers)
	if err != nil {
		return Payload{}, err
	}
	p.SigBundles = sigBundles

	return p, nil
}

// SigningBytes returns a byte slice of version|timestamp|ttl|len|data used as
// the message to be signed by a Signer.
func (p Payload) SigningBytes() []byte {
	j := [][]byte{
		uint64ToBytes(uint64(p.Version)),
		uint64ToBytes(uint64(p.Timestamp.UnixNano())),
		uint64ToBytes(uint64(p.TTL.Nanoseconds())),
		uint64ToBytes(uint64(len(p.Data))),
		p.Data,
	}
	return bytes.Join(j, []byte{})
}

// PubKeyBytes returns a byte slice of all pubkeys concatenated in the index
// order of the slice of sig.Bundles. This is intended to be used with a hash
// function to derive the unique endpoint for a payload on hashmap server.
func (p Payload) PubKeyBytes() []byte {
	return sigutil.BundlePubKeys(p.SigBundles)
}

// PubKeyHash returns a byte slice of the blake2b-512 hash of PubKeyBytes
func (p Payload) PubKeyHash() []byte {
	return sigutil.BundleHash(p.SigBundles)
}

// Endpoint returns a url-safe base64 encoded endpoint string of PubKeyHash
func (p Payload) Endpoint() string {
	return sigutil.EncodedBundleHash(p.SigBundles)
}

// uint64ToBytes converts uint64 numbers into a byte slice in Big Endian format
func uint64ToBytes(t uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, t)
	return b
}
