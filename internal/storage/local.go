package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Local struct{ root string }

func NewLocal(root string) (*Local, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return nil, err
	}
	return &Local{root: filepath.Clean(abs)}, nil
}
func (l *Local) Kind() string                       { return "local" }
func (l *Local) AbsPath(key string) (string, error) { return secureJoin(l.root, key) }
func (l *Local) Put(ctx context.Context, key string, r io.Reader, size int64) error {
	path, err := l.AbsPath(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".nud-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	written, copyErr := io.Copy(tmp, io.LimitReader(r, size+1))
	closeErr := tmp.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if written != size {
		return fmt.Errorf("size mismatch: expected %d, got %d", size, written)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return os.Rename(name, path)
}
func (l *Local) Open(_ context.Context, key string, opt *OpenOptions) (io.ReadCloser, error) {
	path, err := l.AbsPath(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if opt != nil && opt.Offset > 0 {
		if _, err := f.Seek(opt.Offset, io.SeekStart); err != nil {
			_ = f.Close()
			return nil, err
		}
	}
	if opt != nil && opt.Length > 0 {
		return &limitedReadCloser{Reader: io.LimitReader(f, opt.Length), Closer: f}, nil
	}
	return f, nil
}
func (l *Local) Stat(_ context.Context, key string) (*ObjectInfo, error) {
	path, err := l.AbsPath(key)
	if err != nil {
		return nil, err
	}
	s, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{Key: key, Size: s.Size(), ModTime: s.ModTime()}, nil
}
func (l *Local) Delete(_ context.Context, key string) error {
	path, err := l.AbsPath(key)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
func (l *Local) PresignURL(context.Context, string, string, time.Duration) (string, error) {
	return "", ErrPresignUnsupported
}

type limitedReadCloser struct {
	io.Reader
	io.Closer
}

func secureJoin(root, key string) (string, error) {
	if key == "" || filepath.IsAbs(key) || filepath.VolumeName(key) != "" || strings.HasPrefix(key, "/") || strings.HasPrefix(key, `\`) {
		return "", fmt.Errorf("unsafe object key")
	}
	clean := filepath.Clean(filepath.FromSlash(key))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe object key")
	}
	joined := filepath.Join(root, clean)
	rel, err := filepath.Rel(root, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe object key")
	}
	return joined, nil
}
