package cryptoutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadOrCreateKey(dataDir, name, envName string) ([]byte, error) {
	if raw := strings.TrimSpace(os.Getenv(envName)); raw != "" {
		b, err := base64.StdEncoding.DecodeString(raw)
		if err != nil || len(b) != 32 {
			return nil, fmt.Errorf("%s must be base64 for exactly 32 bytes", envName)
		}
		return b, nil
	}
	dir := filepath.Join(dataDir, "secret")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, name)
	if b, err := os.ReadFile(path); err == nil {
		if len(b) != 32 {
			return nil, fmt.Errorf("invalid key length in %s", path)
		}
		return b, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return nil, err
	}
	return b, nil
}
