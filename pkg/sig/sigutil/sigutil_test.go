package sigutil

import (
	"bytes"
	"testing"

	sig "github.com/nomasters/hashmap/pkg/sig"
	// sig "github.com/nomasters/hashmap/pkg/sig"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	s := NewDefaultSigners()
	b, err := Encode(s)
	if err != nil {
		t.Error(err)
	}

	o, err := Decode(b)
	if err != nil {
		t.Error(err)
	}
	m := []byte("hello, world.")
	sb, err := s[0].Sign(m)
	if err != nil {
		t.Error(err)
	}

	ob, err := o[0].Sign(m)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(sb.Sig, ob.Sig) {
		t.Log("source:  ", sb.Sig)
		t.Log("decoded: ", ob.Sig)
		t.Error("Decode signature mismatch")
	}
}

func TestSignAll(t *testing.T) {
	t.Parallel()
	s := NewDefaultSigners()
	m := []byte("hello, world")
	if _, err := SignAll(m, s); err != nil {
		t.Error(err)
	}
}

func TestVerifyAll(t *testing.T) {
	t.Parallel()
	s := NewDefaultSigners()
	m := []byte("hello, world")
	b, err := s[0].Sign(m)
	if err != nil {
		t.Error(err)
	}
	if !VerifyAll(m, []sig.Bundle{b}) {
		t.Error("VerifyAll failed for", b)
	}
}
