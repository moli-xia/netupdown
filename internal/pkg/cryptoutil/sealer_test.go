package cryptoutil

import (
	"bytes"
	"testing"
)

func TestSealerRoundTrip(t *testing.T) {
	s, err := NewSealer(bytes.Repeat([]byte{7}, 32), "storage-config")
	if err != nil {
		t.Fatal(err)
	}
	ciphertext, err := s.Encrypt(`{"secret":"value"}`)
	if err != nil {
		t.Fatal(err)
	}
	if ciphertext == `{"secret":"value"}` {
		t.Fatal("value was not encrypted")
	}
	plain, err := s.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if plain != `{"secret":"value"}` {
		t.Fatalf("got %q", plain)
	}
}
