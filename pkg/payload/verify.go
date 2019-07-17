package payload

import (
	"errors"
	"fmt"
	"time"

	sigutil "github.com/nomasters/hashmap/pkg/sig/sigutil"
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
	// MinTTL is the minimum value of a TTL for a payload
	MinTTL = 0 * time.Second
	// MaxTTL is the maximum value of a TTL for a payload
	MaxTTL = 24 * 7 * time.Hour // 1 week
)

// validateContext is used for interacting with options
type validateContext struct {
	endpoint      string
	ttl           bool
	expiration    bool
	payloadSize   bool
	dataSize      bool
	version       bool
	submitTime    bool
	futureTime    bool
	referenceTime time.Time
}

// WithValidateEndpoint sets the endpoint string for options.validate.endpoint and is
// used for the Verify method. endpoint defaults to and empty string.
func WithValidateEndpoint(e string) Option {
	return func(o *options) {
		o.validate.endpoint = e
	}
}

// WithReferenceTime sets time for options.validate.referenceTime and is used for
// the Verify method. referenceTime defaults to time.Now
func WithReferenceTime(t time.Time) Option {
	return func(o *options) {
		o.validate.referenceTime = t
	}
}

// WithServerMode sets options.validate.submitTime boolean. Defaults to false.
// Setting to false will skip validation when using the payload Verify method
func WithServerMode(b bool) Option {
	return func(o *options) {
		o.validate.submitTime = b
	}
}

// WithValidateTTL sets options.validate.ttl boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateTTL(b bool) Option {
	return func(o *options) {
		o.validate.ttl = b
	}
}

// WithValidateExpiration sets options.validate.expiration boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateExpiration(b bool) Option {
	return func(o *options) {
		o.validate.expiration = b
	}
}

// WithValidateFuture sets options.validate.futureTime boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateFuture(b bool) Option {
	return func(o *options) {
		o.validate.futureTime = b
	}
}

// WithValidateDataSize sets options.validate.dataSize boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateDataSize(b bool) Option {
	return func(o *options) {
		o.validate.dataSize = b
	}
}

// WithValidatePayloadSize sets options.validate.payloadSize boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidatePayloadSize(b bool) Option {
	return func(o *options) {
		o.validate.payloadSize = b
	}
}

// WithValidateVersion sets options.validate.version boolean. Defaults to true.
// Setting to false will skip validation when using the payload Verify method
func WithValidateVersion(b bool) Option {
	return func(o *options) {
		o.validate.version = b
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

	if ok := p.VerifySignatures(); !ok {
		return errors.New("failed signature verification")
	}

	return nil
}

// validate takes a payload and set of options and validates
// the payload itself. This ensures it meets size and version
// requirements
func validate(p Payload, opts ...Option) error {
	o := parseOptions(opts...)

	if o.validate.endpoint != "" {
		if !p.ValidEndpoint(o.validate.endpoint) {
			return errors.New("invalid endpoint")
		}
	}

	if o.validate.payloadSize {
		if !p.ValidPayloadSize() {
			return errors.New("MaxPayloadSize exceeded")
		}
	}
	if o.validate.dataSize {
		if !p.ValidDataSize() {
			return errors.New("MaxMessageSize exceeded")
		}
	}
	if o.validate.version {
		if !p.ValidVersion() {
			return errors.New("invalid payload version")
		}
	}
	if o.validate.expiration {
		if p.IsExpired(o.validate.referenceTime) {
			return errors.New("payload ttl is expired")
		}
	}
	if o.validate.ttl {
		if !p.ValidTTL() {
			return errors.New("invalid payload ttl")
		}
	}
	if o.validate.futureTime {
		if p.IsInFuture(o.validate.referenceTime) {
			return errors.New("payload timestamp is too far in the future")
		}
	}
	if o.validate.submitTime {
		if !p.WithinSubmitWindow(o.validate.referenceTime) {
			return errors.New("timestamp is outside of submit window")
		}
	}
	return nil
}

// ValidEndpoint takes a string and attempts to match the URL safe
// base64 string encoded PubKeyHash and returns a boolean
func (p Payload) ValidEndpoint(e string) bool {
	return e == p.Endpoint()
}

// ValidTTL checks that a TTL falls within an acceptable range.
func (p Payload) ValidTTL() bool {
	return !(p.TTL < MinTTL || p.TTL > MaxTTL)
}

// IsExpired checks the reference time t against the timestamp and
// TTL of a payload and returns a boolean value on whether
// or not the TTL has been exceeded
func (p Payload) IsExpired(t time.Time) bool {
	return t.Sub(p.Timestamp) > p.TTL
}

// IsInFuture checks if the payload timestamp is too far into the future based
// on the reference time t plus the MaxSubmitWindow.
func (p Payload) IsInFuture(t time.Time) bool {
	return p.Timestamp.UnixNano() > t.Add(MaxSubmitWindow).UnixNano()
}

// ValidVersion returns whether version is supported by Hashmap
// Currently only V1 is supported.
func (p Payload) ValidVersion() bool {
	switch p.Version {
	case V1:
		return true
	}
	return false
}

// ValidDataSize checks that the length of Payload.Data is less than or equal
// to the MaxMessageSize and returns a boolean value.
func (p Payload) ValidDataSize() bool {
	return len(p.Data) <= MaxMessageSize
}

// ValidPayloadSize checks that the wire protocol bytes are less than or equal
// to the MaxPayloadSize allowed and returns a boolean value.
func (p Payload) ValidPayloadSize() bool {
	b, err := Marshal(p)
	if err != nil {
		return false
	}
	if len(b) > MaxPayloadSize {
		return false
	}
	return true
}

// WithinSubmitWindow checks reference time t against the payload timestamp,
// validates that it exists within the MaxSubmitWindow and returns a boolean.
func (p Payload) WithinSubmitWindow(t time.Time) bool {
	diff := t.Sub(p.Timestamp)

	// get absolute value of time difference
	if diff.Seconds() < 0 {
		diff = -diff
	}
	return diff <= MaxSubmitWindow
}

// VerifySignatures checks all signatures in the sigBundles. If all signatures
// are valid, it returns `true`.
func (p Payload) VerifySignatures() bool {
	return sigutil.VerifyAll(p.SigningBytes(), p.SigBundles)
}
