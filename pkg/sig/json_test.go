package sig

import (
	"encoding/json"
	"testing"
)

func TestJSONEncoding(t *testing.T) {
	t.Parallel()
	s := GenNaclSign()
	m := []byte("sign me, plz.")
	b, _ := s.Sign(m)
	if _, err := json.Marshal(b); err != nil {
		t.Error(err)
	}
}

func TestUnmarshall(t *testing.T) {
	t.Parallel()
	valid := `
		{
			"alg":1,
			"pub":"qOQp8p/ZheaTHhlJ90TQoHBNKQ7BXGC2FHOjNHfGH0M=",
			"sig":"4xb4Zgr9lGN3imfiGezH36fby6cktqrpRL6m5yluqSn83r+HMONhQ5722BDgF5Cb4xX9GVUiWo6I/zFOeGlTBQ=="
		}`
	invalidString := `
		{
			"alg":1,
			"pub":"invalid",
			"sig":"4xb4Zgr9lGN3imfiGezH36fby6cktqrpRL6m5yluqSn83r+HMONhQ5722BDgF5Cb4xX9GVUiWo6I/zFOeGlTBQ=="
		}`
	invalidType := `
		{
			"alg":1,
			"pub":1,
			"sig":"4xb4Zgr9lGN3imfiGezH36fby6cktqrpRL6m5yluqSn83r+HMONhQ5722BDgF5Cb4xX9GVUiWo6I/zFOeGlTBQ=="
		}`

	var b Bundle
	if err := json.Unmarshal([]byte(valid), &b); err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal([]byte(invalidString), &b); err == nil {
		t.Error("invalid pub encoding not caught")
	}
	if err := json.Unmarshal([]byte(invalidType), &b); err == nil {
		t.Error("invalid type encoding not caught")
	}
}
