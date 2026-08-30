package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Sealer encrypts sensitive database values with an independently-derived key.
type Sealer struct{ aead cipher.AEAD }

func NewSealer(root []byte, info string) (*Sealer, error) {
	if len(root) != 32 {
		return nil, errors.New("root encryption key must be 32 bytes")
	}
	h := hmac.New(sha256.New, root)
	_, _ = h.Write([]byte("netupdown:" + info))
	block, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Sealer{aead: aead}, nil
}

func (s *Sealer) Encrypt(plain string) (string, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := s.aead.Seal(nil, nonce, []byte(plain), []byte("netupdown:v1"))
	return "enc:v1:" + base64.RawStdEncoding.EncodeToString(nonce) + ":" + base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

func (s *Sealer) Decrypt(value string) (string, error) {
	if !strings.HasPrefix(value, "enc:v1:") {
		return value, nil
	}
	parts := strings.Split(value, ":")
	if len(parts) != 4 {
		return "", errors.New("invalid encrypted value format")
	}
	nonce, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("decode nonce: %w", err)
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	plain, err := s.aead.Open(nil, nonce, ciphertext, []byte("netupdown:v1"))
	if err != nil {
		return "", fmt.Errorf("decrypt value: %w", err)
	}
	return string(plain), nil
}
