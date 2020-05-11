package analyze

import (
	"time"

	payload "github.com/nomasters/hashmap/pkg/payload"
)

// Payload is the analyze struct container for the standard
// output.
type Payload struct {
	Raw                payload.Payload `json:"raw_payload"`
	Hash               string          `json:"endpoint_hash"`
	Timestamp          time.Time       `json:"timestamp"`
	TTL                string          `json:"ttl"`
	Expired            bool            `json:"expired"`
	ValidVersion       bool            `json:"valid_version"`
	ValidTimestamp     bool            `json:"valid_timestamp"`
	WithinSubmitWindow bool            `json:"within_submit_window"`
	ValidTTL           bool            `json:"valid_ttl"`
	ValidSignatures    bool            `json:"valid_signatures"`
	ValidDataSize      bool            `json:"valid_data_size"`
	ValidPayloadSize   bool            `json:"valid_payload_size"`
	ErrorMessage       string          `json:"error_message,omitempty"`
}

// NewPayload returns a payload analysis and runs the entire validation suite on the output
func NewPayload(b []byte) (*Payload, error) {
	var p Payload
	pl, err := payload.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	p.Raw = pl
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
	p.ValidDataSize = pl.ValidDataSize()
	p.ValidPayloadSize = pl.ValidPayloadSize()
}
