package payload

import (
	"errors"
	"fmt"
	"time"

	"github.com/nomasters/hashmap/sig"
)

const (
	// MaxSigBundleCount is the upper limit of signatures
	// allowed for a single payload
	MaxSigBundleCount = 4
	// MaxPayloadSize is intended to by used by a LimitedReader
	// to enforce a strict upper limit on payload size
	MaxPayloadSize = 128 * 1024 // 128 KB
	// MaxMessageSize Payload.Data
	MaxMessageSize = 512 // bytes
	// MaxSubmitWindow is the time drift allow between a submission
	// to hashmap server and the time reflected on a signed payload
	MaxSubmitWindow = 5 * time.Second
)

// validateContext is used for interacting with Context
type validateContext struct {
	ttl           bool
	payloadSize   bool
	dataSize      bool
	version       bool
	submitTime    bool
	futureTime    bool
	referenceTime time.Time
}

// WithReferenceTime sets time for Context.validate.referenceTime and is used for
// the Verify method. referenceTime defaults to time.Now
func WithReferenceTime(t time.Time) Option {
	return func(c *Context) {
		c.validate.referenceTime = t
	}
}

// WithServerMode sets context.validate.submitTime boolean. Defaults to false.
// Setting to false will skip validation when using the payload Verify method
func WithServerMode(b bool) Option {
	return func(c *Context) {
		c.validate.submitTime = b
	}
}

// WithValidateTTL sets context.validate.ttl boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateTTL(b bool) Option {
	return func(c *Context) {
		c.validate.ttl = b
	}
}

// WithValidateFuture sets context.validate.futureTime boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateFuture(b bool) Option {
	return func(c *Context) {
		c.validate.futureTime = b
	}
}

// WithValidateDataSize sets context.validate.dataSize boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateDataSize(b bool) Option {
	return func(c *Context) {
		c.validate.dataSize = b
	}
}

// WithValidatePayloadSize sets context.validate.payloadSize boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidatePayloadSize(b bool) Option {
	return func(c *Context) {
		c.validate.payloadSize = b
	}
}

// WithValidateVersion sets context.validate.version boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateVersion(b bool) Option {
	return func(c *Context) {
		c.validate.version = b
	}
}

// Verify method takes a set of options and implements the Verify function
func (p Payload) Verify(options ...Option) error {
	return verify(p, options...)
}

// verify takes a payload and set of options and validates
// and verifies the payload. By default Verify runs in client mode
// which means Verification passes without verifying the submitWindow.
// Any host planning to store hashmap payloads should run WithServerMode
func verify(p Payload, options ...Option) error {

	if err := validate(p, options...); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if ok := p.verifySignatures(); !ok {
		return errors.New("failed signature verification")
	}

	return nil
}

// validate takes a payload and set of options and validates
// the payload itself. This ensures it meets size and version
// requirements
func validate(p Payload, options ...Option) error {
	c := parseOptions(options...)
	if c.validate.payloadSize {
		if !p.validPayloadSize() {
			return errors.New("MaxPayloadSize exceeded")
		}
	}
	if c.validate.dataSize {
		if !p.validDataSize() {
			return errors.New("MaxMessageSize exceeded")
		}
	}
	if c.validate.version {
		if !p.validVersion() {
			return errors.New("invalid payload version")
		}
	}
	if c.validate.ttl {
		if p.isExpired(c.validate.referenceTime) {
			return errors.New("payload ttl is expired")
		}
	}
	if c.validate.futureTime {
		if p.isInFuture(c.validate.referenceTime) {
			return errors.New("payload timestamp is too far in the future")
		}
	}
	if c.validate.submitTime {
		if !p.withinSubmitWindow(c.validate.referenceTime) {
			return errors.New("timestamp is outside of submit window")
		}
	}
	return nil
}

// isExpired checks the reference time t against the timestamp and
// TTL of a payload and returns a boolean value on whether
// or not the TTL has been exceeded
func (p Payload) isExpired(t time.Time) bool {
	return t.Sub(p.Timestamp) > p.TTL
}

// isInFuture checks if the payload timestamp is too far into the future based
// on the reference time t plus the MaxSubmitWindow.
func (p Payload) isInFuture(t time.Time) bool {
	return p.Timestamp.UnixNano() > t.Add(MaxSubmitWindow).UnixNano()
}

// validVersion returns whether version is supported by Hashmap
// Currently only V1 is supported.
func (p Payload) validVersion() bool {
	switch p.Version {
	case V1:
		return true
	}
	return false
}

// validDataSize checks that the length of Payload.Data is less than or equal
// to the MaxMessageSize and returns a boolean value.
func (p Payload) validDataSize() bool {
	return len(p.Data) <= MaxMessageSize
}

// validPayloadSize checks that the wire protocol bytes are less than or equal
// to the MaxPayloadSize allowed and returns a boolean value.
func (p Payload) validPayloadSize() bool {
	b, err := Marshal(p)
	if err != nil {
		return false
	}
	if len(b) > MaxPayloadSize {
		return false
	}
	return true
}

// withinSubmitWindow checks reference time t against the payload timestamp,
// validates that it exists within the MaxSubmitWindow and returns a boolean.
func (p Payload) withinSubmitWindow(t time.Time) bool {
	diff := t.Sub(p.Timestamp)

	// get absolute value of time difference
	if diff.Seconds() < 0 {
		diff = -diff
	}
	return diff <= MaxSubmitWindow
}

// verifySignatures checks all signatures in the sigBundles. If all signatures
// are valid, it returns `true`.
func (p Payload) verifySignatures() bool {
	for _, bundle := range p.SigBundles {
		if ok := sig.Verify(p.SigningBytes(), bundle); !ok {
			return false
		}
	}
	return true
}
