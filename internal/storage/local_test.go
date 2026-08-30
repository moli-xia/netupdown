package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestSecureJoin(t *testing.T) {
	root := t.TempDir()
	for _, key := range []string{"../x", "../../etc/passwd", `C:\x`, `\\server\share`, `/absolute`} {
		if _, err := secureJoin(root, key); err == nil {
			t.Fatalf("expected %q rejected", key)
		}
	}
	if _, err := secureJoin(root, "app/1/file.zip"); err != nil {
		t.Fatal(err)
	}
}
func TestLocalRoundTrip(t *testing.T) {
	d, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("0123456789")
	if err := d.Put(context.Background(), "a/b.bin", bytes.NewReader(data), int64(len(data))); err != nil {
		t.Fatal(err)
	}
	r, err := d.Open(context.Background(), "a/b.bin", &OpenOptions{Offset: 2, Length: 4})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	got, _ := io.ReadAll(r)
	if string(got) != "2345" {
		t.Fatalf("got %q", got)
	}
}
