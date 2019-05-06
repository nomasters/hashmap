//go:generate protoc -I=pb --go_out=pb/ pb/payload.proto

package payload

import (
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

// Verifier is an interface for verifying Payloads and can take an arbitrary number of Options
// to modify the strictness of the verification.
type Verifier interface {
	Verify(options ...Option) error
}

// Option is used for interacting with Context for Options on Generate and
// Verified tooling.
type Option func(*Context)

// Context contains private fields used for Option
type Context struct {
	version         Version
	timestamp       time.Time
	ttl             time.Duration
	verifyTTL       bool
	verifyTimestamp bool
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

// TODO:
// - Sort out what Options we should pass in: I'm thinking TTL, Timestamp, and Version
// -

func Generate(message []byte, signers []sig.Signer, options ...Option) (Payload, error) {
	return Payload{}, nil
}
