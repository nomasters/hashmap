package analyze

import (
	"time"
)

type PayloadAnalysis struct {
	Hash       string
	Age        time.Duration
	ExpiresIn  time.Duration
	Payload    Payload
	Validation Validation
}

type Payload struct {
	Version    string
	Timestamp  time.Time
	TTL        time.Duration
	SigBundles []Bundle
	Len        int
	Data       string
}

type Bundle struct {
	Alg string
	Pub string
	Sig string
}

type Validation struct {
	Version        bool
	Timestamp      bool
	InSubmitWindow bool
	TTL            bool
	Expired        bool
	Signatures     []bool
	Length         bool
	DataSize       bool
	Errors         []string
}