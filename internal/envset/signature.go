package envset

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

var (
	ErrSignatureMismatch = errors.New("signature mismatch: envset may have been tampered with")
	ErrSignatureNotFound = errors.New("signature not found on envset")
	ErrEmptyPassphrase   = errors.New("passphrase must not be empty")
)

const signatureMetaKey = "__signature__"

// Sign computes an HMAC-SHA256 signature over all key-value pairs in the
// EnvSet (sorted by key) and stores it as metadata. Existing signatures
// are replaced.
func Sign(es *EnvSet, passphrase string) error {
	if es == nil {
		return errors.New("envset must not be nil")
	}
	if passphrase == "" {
		return ErrEmptyPassphrase
	}

	sig, err := computeSignature(es, passphrase)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	if es.Meta == nil {
		es.Meta = make(map[string]string)
	}
	es.Meta[signatureMetaKey] = sig
	return nil
}

// VerifySignature checks the stored HMAC-SHA256 signature against the
// current key-value contents of the EnvSet.
func VerifySignature(es *EnvSet, passphrase string) error {
	if es == nil {
		return errors.New("envset must not be nil")
	}
	if passphrase == "" {
		return ErrEmptyPassphrase
	}

	stored, ok := es.Meta[signatureMetaKey]
	if !ok || stored == "" {
		return ErrSignatureNotFound
	}

	expected, err := computeSignature(es, passphrase)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	if !hmac.Equal([]byte(stored), []byte(expected)) {
		return ErrSignatureMismatch
	}
	return nil
}

// ClearSignature removes the stored signature from the EnvSet metadata.
func ClearSignature(es *EnvSet) {
	if es == nil || es.Meta == nil {
		return
	}
	delete(es.Meta, signatureMetaKey)
}

// HasSignature reports whether the EnvSet carries a stored signature.
func HasSignature(es *EnvSet) bool {
	if es == nil || es.Meta == nil {
		return false
	}
	_, ok := es.Meta[signatureMetaKey]
	return ok
}

func computeSignature(es *EnvSet, passphrase string) (string, error) {
	keys := make([]string, 0, len(es.Vars))
	for k := range es.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	mac := hmac.New(sha256.New, []byte(passphrase))
	for _, k := range keys {
		if k == signatureMetaKey {
			continue
		}
		_, err := fmt.Fprintf(mac, "%s=%s\n", k, es.Vars[k])
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
