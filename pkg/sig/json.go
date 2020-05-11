package sig

import (
	"encoding/base64"
	"encoding/json"
)

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

// UnmarshalJSON unmarshals base64 encoded strings into Bytes
func (b Bytes) MarshalJSON() ([]byte, error) {
	s := base64.StdEncoding.EncodeToString(b)
	return json.Marshal(s)
}
