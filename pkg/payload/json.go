package payload

import (
	"encoding/base64"
	"encoding/json"
)

// Unmarshal takes encoded json and returns a Payload and error
func Unmarshal(b []byte) (Payload, error) {
	var p Payload
	err := json.Unmarshal(b, &p)
	return p, err
}

// Marshal takes a payload and returns encoded JSON and error
func Marshal(p Payload) ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalJSON unmarshals base64 encoded strings into Bytes
func (b *Bytes) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	*b = data
	return nil
}

// MarshalJSON takes Bytes and returns a base64 encoded string
func (b Bytes) MarshalJSON() ([]byte, error) {
	s := base64.StdEncoding.EncodeToString(b)
	return json.Marshal(s)
}
