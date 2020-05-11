package analyze

import (
	"io/ioutil"
	"testing"
)

func TestNewPayload(t *testing.T) {
	validExpired := "../../test/testdata/valid_payload_expired.json"
	protoBytes, err := ioutil.ReadFile(validExpired)
	if err != nil {
		t.Error(err)
	}
	if _, err := NewPayload(protoBytes); err != nil {
		t.Error(err)
	}
}
