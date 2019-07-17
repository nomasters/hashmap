package analyze

import (
	"time"

	proto "github.com/golang/protobuf/proto"
	payload "github.com/nomasters/hashmap/pkg/payload"
	pb "github.com/nomasters/hashmap/pkg/payload/pb"
)

// Payload is the analyze struct container for the standard
// output.
type Payload struct {
	Raw                pb.Payload `json:"raw_payload"`
	Hash               string     `json:"endpoint_hash"`
	Timestamp          time.Time  `json:"timestamp"`
	TTL                string     `json:"ttl"`
	Expired            bool       `json:"expired"`
	ValidVersion       bool       `json:"valid_version"`
	ValidTimestamp     bool       `json:"valid_timestamp"`
	WithinSubmitWindow bool       `json:"within_submit_window"`
	ValidTTL           bool       `json:"valid_ttl"`
	ValidSignatures    bool       `json:"valid_signatures"`
	ValidLength        bool       `json:"valid_length"`
	ValidDataSize      bool       `json:"valid_data_size"`
	ValidPayloadSize   bool       `json:"valid_payload_size"`
	ErrorMessage       string     `json:"error_message,omitempty"`
}

// NewPayload returns a payload analysis and runs the entire validation suite on the output
func NewPayload(b []byte) (*Payload, error) {
	var p Payload
	if err := proto.Unmarshal(b, &p.Raw); err != nil {
		return nil, err
	}
	pl, err := payload.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	p.Hash = pl.Endpoint()
	p.Timestamp = pl.Timestamp
	p.TTL = pl.TTL.String()
	p.analyze(pl)
	if err := pl.Verify(); err != nil {
		p.ErrorMessage = err.Error()
	}
	return &p, nil
}

func (p *Payload) analyze(pl payload.Payload) {
	now := time.Now()
	p.Expired = pl.IsExpired(now)
	p.ValidVersion = pl.ValidVersion()
	p.ValidTimestamp = !pl.IsInFuture(now)
	p.WithinSubmitWindow = pl.WithinSubmitWindow(now)
	p.ValidTTL = pl.ValidTTL()
	p.ValidSignatures = pl.VerifySignatures()
	p.ValidLength = int(p.Raw.Len) == len(p.Raw.Data)
	p.ValidDataSize = pl.ValidDataSize()
	p.ValidPayloadSize = pl.ValidPayloadSize()
}
