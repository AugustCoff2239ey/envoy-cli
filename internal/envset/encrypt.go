package envset

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the encryption key is empty or invalid.
var ErrInvalidKey = errors.New("encrypt: key must not be empty")

// ErrDecryptFailed is returned when decryption fails (e.g. wrong key or corrupt data).
var ErrDecryptFailed = errors.New("encrypt: decryption failed, invalid key or corrupted data")

// deriveKey derives a 32-byte AES-256 key from an arbitrary passphrase.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// Encrypt encrypts the plaintext using AES-256-GCM with the given passphrase.
// The returned string is base64-encoded and includes the nonce prepended.
func Encrypt(plaintext, passphrase string) (string, error) {
	if passphrase == "" {
		return "", ErrInvalidKey
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext produced by Encrypt.
func Decrypt(encoded, passphrase string) (string, error) {
	if passphrase == "" {
		return "", ErrInvalidKey
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrDecryptFailed
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", ErrDecryptFailed
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}
	return string(plaintext), nil
}
