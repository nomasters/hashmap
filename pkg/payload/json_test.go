package payload

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	sig "github.com/nomasters/hashmap/pkg/sig"
)

func TestJSONMarshall(t *testing.T) {
	t.Parallel()

	m := []byte("sign me, plz.")
	s := []sig.Signer{sig.GenNaclSign()}
	d, err := Generate(m, s, WithTTL(5*time.Second))
	if err != nil {
		t.Error(err)
	}
	blob, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(blob))
}

func TestJSONUnmarshall(t *testing.T) {
	t.Parallel()
	valid := `
		{
			"version":1,
			"timestamp":"2020-05-11T05:56:43.839934-06:00",
			"ttl":86400000000000,
			"sig_bundles":[
				{
					"alg":1,
					"pub":"jpZdClX/d44twVle7gbVc61Ai4Tv8xSvHRAlhagbWq4=",
					"sig":"zi5awXZr85U6ecR/ffAuHFzeNmLWEkOW0as6GMJqhYLbLDn+dovqGJSoeJmbFJEUqK4+Q6gMPrXmi+ngXz/mCw=="
				}
			],
			"data":"c2lnbiBtZSwgcGx6Lg=="
		}`
	invalidString := `
		{
			"version":1,
			"timestamp":"2020-05-11T05:56:43.839934-06:00",
			"ttl":86400000000000,
			"sig_bundles":[
				{
					"alg":1,
					"pub":"jpZdClX/d44twVle7gbVc61Ai4Tv8xSvHRAlhagbWq4=",
					"sig":"zi5awXZr85U6ecR/ffAuHFzeNmLWEkOW0as6GMJqhYLbLDn+dovqGJSoeJmbFJEUqK4+Q6gMPrXmi+ngXz/mCw=="
				}
			],
			"data":"bad_string"
		}`
	invalidType := `
		{
			"version":1,
			"timestamp":"2020-05-11T05:56:43.839934-06:00",
			"ttl":86400000000000,
			"sig_bundles":[
				{
					"alg":1,
					"pub":"jpZdClX/d44twVle7gbVc61Ai4Tv8xSvHRAlhagbWq4=",
					"sig":"zi5awXZr85U6ecR/ffAuHFzeNmLWEkOW0as6GMJqhYLbLDn+dovqGJSoeJmbFJEUqK4+Q6gMPrXmi+ngXz/mCw=="
				}
			],
			"data":1
		}`

	var p Payload
	if err := json.Unmarshal([]byte(valid), &p); err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal([]byte(invalidString), &p); err == nil {
		t.Error("invalid pub encoding not caught")
	}
	if err := json.Unmarshal([]byte(invalidType), &p); err == nil {
		t.Error("invalid type encoding not caught")
	}
	t.Log(p.TTL)
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	var signers []sig.Signer
	signers = append(signers, sig.GenNaclSign())
	message := []byte("hello, world")

	t.Run("Normal Operation", func(t *testing.T) {
		p, err := Generate(message, signers)
		if err != nil {
			t.Error(err)
		}
		if _, err := Marshal(p); err != nil {
			t.Error(err)
		}
	})

	t.Run("Invalid Timestamp", func(t *testing.T) {
		bad := time.Unix(-99999999999, 0)
		p, err := Generate(message, signers, WithTimestamp(bad))
		if err != nil {
			t.Error(err)
		}
		if _, err := Marshal(p); err == nil {
			t.Error("failed to catch invalid timestamp")
		}
	})

}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	t.Run("Normal Operation", func(t *testing.T) {
		var signers []sig.Signer
		signers = append(signers, sig.GenNaclSign())
		message := []byte("hello, world")
		p, err := Generate(message, signers)
		if err != nil {
			t.Error(err)
		}
		encoded, err := Marshal(p)
		if err != nil {
			t.Error(err)
		}
		if _, err := Unmarshal(encoded); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid protobuf", func(t *testing.T) {
		invalidPayload, _ := hex.DecodeString("ffffffffffffff")
		// this is a valid payload protobuf with a malformed (by hand) timestamp to force out of range errors
		badTimestamp, _ := hex.DecodeString("12070880cdfdffaf071a080880aef188221001")
		// this is a valid payload protobuf with a malformed (by hand) ttl to force out of range errors
		badTTL, _ := hex.DecodeString("12070880cdfdefaf071a0808fffffff8221001")

		testTable := []struct {
			bytes []byte
			err   string
		}{
			{
				bytes: invalidPayload,
				err:   "failed to catch malformed payload bytes",
			},
			{
				bytes: badTimestamp,
				err:   "failed to catch malformed timestamp",
			},
			{
				bytes: badTTL,
				err:   "failed to catch malformed ttl",
			},
		}

		for _, test := range testTable {
			if _, err := Unmarshal(test.bytes); err == nil {
				t.Error(test.err)
			}
		}
	})
}
