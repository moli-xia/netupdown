package storage

import "testing"

func TestS3ObjectKey(t *testing.T) {
	d := &S3{cfg: S3Config{BasePath: "files/"}}
	if got := d.objectKey("app/1/file.zip"); got != "files/app/1/file.zip" {
		t.Fatalf("got %q", got)
	}
}
