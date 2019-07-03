package analyze

import (
	"time"
)

type Analysis struct {
	Metadata   Metadata
	Payload    Payload
	Validation Validation
}

type Metadata struct {
	Hash      string
	Age       time.Duration
	ExpiresIn time.Duration
}

type Payload struct {
	Version    string
	Timestamp  time.Time
	TTL        time.Duration
	SigBundles []Bundle
	Data       string
}

type Bundle struct {
	Alg string
	Pub string
	Sig string
}

type Validation struct {
	Version         bool
	TTL             bool
	Timestamp       bool
	InSubmitWindow  bool
	DataSize        bool
	Signatures      []bool
	FailureMessages []string
}
