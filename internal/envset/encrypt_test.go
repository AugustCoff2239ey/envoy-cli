package envset

import (
	"strings"
	"testing"
)

func TestEncrypt_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	passphrase := "my-strong-passphrase"

	encoded, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt returned unexpected error: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encrypt returned empty string")
	}

	decoded, err := Decrypt(encoded, passphrase)
	if err != nil {
		t.Fatalf("Decrypt returned unexpected error: %v", err)
	}
	if decoded != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decoded)
	}
}

func TestEncrypt_UniqueOutputs(t *testing.T) {
	plaintext := "hello"
	passphrase := "key"

	a, _ := Encrypt(plaintext, passphrase)
	b, _ := Encrypt(plaintext, passphrase)

	// Nonce is random, so outputs should differ
	if a == b {
		t.Error("expected different ciphertexts for the same input due to random nonce")
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := Encrypt("value", "")
	if err != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	_, err := Decrypt("somedata", "")
	if err != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encoded, err := Encrypt("secret", "correct-key")
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	_, err = Decrypt(encoded, "wrong-key")
	if err != ErrDecryptFailed {
		t.Errorf("expected ErrDecryptFailed, got %v", err)
	}
}

func TestDecrypt_CorruptData(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", "key")
	if err != ErrDecryptFailed {
		t.Errorf("expected ErrDecryptFailed, got %v", err)
	}
}

func TestEncrypt_LargeValue(t *testing.T) {
	plaintext := strings.Repeat("A", 10_000)
	passphrase := "large-value-key"

	encoded, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	decoded, err := Decrypt(encoded, passphrase)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if decoded != plaintext {
		t.Error("round-trip failed for large value")
	}
}
