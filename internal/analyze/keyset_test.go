package analyze

import (
	"encoding/json"

	"io/ioutil"
	"testing"
)

func TestNewKeySet(t *testing.T) {
	keysetPath := "../../test/testdata/hashmap_ed25519.keyset"
	b, err := ioutil.ReadFile(keysetPath)
	if err != nil {
		t.Error(err)
	}
	k, err := NewKeySet(b)
	if err != nil {
		t.Error(err)
	}

	output, err := json.Marshal(*k)
	t.Log(string(output))

}
