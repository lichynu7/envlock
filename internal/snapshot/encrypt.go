package snapshot

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptOptions controls encryption behaviour.
type EncryptOptions struct {
	Passphrase string
}

// DefaultEncryptOptions returns sensible defaults (no passphrase).
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{}
}

// deriveKey produces a 32-byte AES-256 key from an arbitrary passphrase.
func deriveKey(passphrase string) []byte {
	sum := sha256.Sum256([]byte(passphrase))
	return sum[:]
}

// Encrypt encodes plaintext using AES-256-GCM with the supplied passphrase
// and returns a base64-encoded ciphertext string.
func Encrypt(plaintext []byte, opts EncryptOptions) (string, error) {
	if opts.Passphrase == "" {
		return "", errors.New("encrypt: passphrase must not be empty")
	}
	block, err := aes.NewCipher(deriveKey(opts.Passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt reverses Encrypt, returning the original plaintext.
func Decrypt(encoded string, opts EncryptOptions) ([]byte, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(deriveKey(opts.Passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("decrypt: ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
