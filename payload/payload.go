//go:generate protoc -I=pb --go_out=pb/ pb/payload.proto

package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	proto "github.com/gogo/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"
	pb "github.com/nomasters/hashmap/payload/pb"
	sig "github.com/nomasters/hashmap/sig"
)

// Version type is used for setting the hashmap implementation version.
type Version int32

const (
	// V0 is deprecated and should be deemed invalid
	V0 Version = iota
	// V1 is the current version of the payload spec
	V1
)

const (
	defaultTTL = 24 * time.Hour
)

var (
	defaultVersion = V1
)

// Payload holds all information related to a Hashmap Payload that will be handled
// for signing and validation. This struct is used by both client and server and
// includes all necessary methods for encoding, decoding, signing, an verifying itself.
type Payload struct {
	Version    Version
	Timestamp  time.Time
	TTL        time.Duration
	SigBundles []sig.Bundle
	Data       []byte
}

// Option is used for interacting with Context when setting options for Generate and Verify
type Option func(*Context)

// Context contains private fields used for Option
type Context struct {
	version   Version
	timestamp time.Time
	ttl       time.Duration
	validate  validateContext
}

// WithVersion takes a Version and returns an Option
func WithVersion(v Version) Option {
	return func(c *Context) {
		c.version = v
	}
}

// WithTimestamp takes a time.Time and returns an Option
func WithTimestamp(t time.Time) Option {
	return func(c *Context) {
		c.timestamp = t
	}
}

// WithTTL takes a time.Duration and returns an Option
func WithTTL(d time.Duration) Option {
	return func(c *Context) {
		c.ttl = d
	}
}

// parseOptions takes a arbitrary number of Option funcs and returns a Context with defaults
// for version, timestamp, and ttl, and validate rules.
func parseOptions(options ...Option) Context {
	now := time.Now()

	c := Context{
		version:   defaultVersion,
		timestamp: now,
		ttl:       defaultTTL,
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
	for _, option := range options {
		option(&c)
	}
	return c
}

// Unmarshal takes a byte slice and attempts to decode the protobuf wire
// format into a Payload. This does not apply any Payload verification.
// Verification should happen after Unmarshalling.
func Unmarshal(b []byte) (Payload, error) {
	var pbp pb.Payload
	if err := proto.Unmarshal(b, &pbp); err != nil {
		return Payload{}, err
	}

	timestamp, err := ptypes.Timestamp(pbp.Timestamp)
	if err != nil {
		return Payload{}, err
	}
	ttl, err := ptypes.Duration(pbp.Ttl)
	if err != nil {
		return Payload{}, err
	}

	var sigBundles []sig.Bundle
	for _, sigSet := range pbp.SigBundles {
		s := sig.Bundle{
			Alg: sig.Alg(sigSet.Alg),
			Pub: sigSet.Pub,
			Sig: sigSet.Sig,
		}
		sigBundles = append(sigBundles, s)
	}

	p := Payload{
		Version:    Version(pbp.Version),
		Timestamp:  timestamp,
		TTL:        ttl,
		SigBundles: sigBundles,
		Data:       pbp.Data,
	}

	return p, nil
}

// Marshal takes a Payload and encodes it into the protobuf wire format.
// This does not apply any Payload verification. Verification should happen
// before marshalling.
func Marshal(p Payload) ([]byte, error) {

	timestamp, err := ptypes.TimestampProto(p.Timestamp)
	if err != nil {
		return []byte{}, err
	}

	var sigBundles []*pb.Payload_SigBundle
	for _, sigBundle := range p.SigBundles {
		s := &pb.Payload_SigBundle{
			Alg: pb.Payload_SigBundle_Alg(sigBundle.Alg),
			Pub: sigBundle.Pub,
			Sig: sigBundle.Sig,
		}
		sigBundles = append(sigBundles, s)
	}

	pbp := &pb.Payload{
		Version:    pb.Payload_Version(p.Version),
		Timestamp:  timestamp,
		Ttl:        ptypes.DurationProto(p.TTL),
		SigBundles: sigBundles,
		Len:        uint32(len(p.Data)),
		Data:       p.Data,
	}

	return proto.Marshal(pbp)

}

// Generate takes a message, signers, and a set of options and returns a payload or error.
// This function defaults to time.Now() and the default TTL of 24 hours. Generate Requires
// at least one signer, but can sign with many signers. Sort order is important though, The unique
// order of the signers pubkeys are what is responsible for generating the endpoint hash.
func Generate(message []byte, signers []sig.Signer, options ...Option) (Payload, error) {
	if len(signers) == 0 {
		return Payload{}, errors.New("Generate must have at least one signer")
	}
	c := parseOptions(options...)
	p := Payload{
		Version:   c.version,
		Timestamp: c.timestamp,
		TTL:       c.ttl,
		Data:      message,
	}

	var sigBundles []sig.Bundle
	for _, s := range signers {
		b, err := s.Sign(p.SigningBytes())
		if err != nil {
			return Payload{}, err
		}
		sigBundles = append(sigBundles, b)
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
	var o []byte
	for _, b := range p.SigBundles {
		o = append(o, b.Pub...)
	}
	return o
}

// uint64ToBytes converts uint64 numbers into a byte slice in Big Endian format
func uint64ToBytes(t uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, t)
	return b
}
