package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrPresignUnsupported = errors.New("storage: presign unsupported")

type ObjectInfo struct {
	Key     string
	Size    int64
	ModTime time.Time
}
type OpenOptions struct {
	Offset int64
	Length int64
}
type Driver interface {
	Kind() string
	Put(context.Context, string, io.Reader, int64) error
	Open(context.Context, string, *OpenOptions) (io.ReadCloser, error)
	Stat(context.Context, string) (*ObjectInfo, error)
	Delete(context.Context, string) error
	PresignURL(context.Context, string, string, time.Duration) (string, error)
}
