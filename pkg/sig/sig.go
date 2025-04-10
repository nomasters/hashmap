package sig

// Alg type is used for setting the Algorithm in a Signature Set
type Alg uint16

// Bytes is a byte slice with special json encoding properties
type Bytes []byte

const (
	// 0 value skipped to handle empty Payload validation
	_ Alg = iota
	// AlgNaClSign is meant for Nacl Sign implementations
	AlgNaClSign
	// AlgXMSS10 is meant for xmss sha2_10_256
	AlgXMSS10
)

// Bundle is used to encapsulate an Algorithm implementation, A Public Key, and a Signature.
// A Bundle is designed to be used to verify the integrity of the Payload.
type Bundle struct {
	Alg Alg   `json:"alg"`
	Pub Bytes `json:"pub"`
	Sig Bytes `json:"sig"`
}

// Signer is an interface for signing messages and generating a SigSet.
type Signer interface {
	Sign(message []byte) (Bundle, error)
}

// Verify takes a message and a signature Bundle and attempts to verify
// the bundle based on bundle's implemented Alg, the sig, and the pubkey.
// Verify returns a simple true or false.
func Verify(message []byte, bundle Bundle) bool {
	switch bundle.Alg {
	case AlgNaClSign:
		return VerifyNaclSign(message, bundle)
	case AlgXMSS10:
		return VerifyXMSS10(message, bundle)
	}
	return false
}
