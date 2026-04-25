package envset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// DigestResult holds the computed digest and per-key contributions.
type DigestResult struct {
	Digest  string            `json:"digest"`
	Entries map[string]string `json:"entries"`
}

// Digest computes a deterministic SHA-256 digest over the key-value pairs of
// an EnvSet. Keys are sorted before hashing so the result is stable regardless
// of insertion order. An optional list of keys can be provided to restrict
// which entries are included; if empty all keys are used.
func Digest(es *EnvSet, keys []string) (*DigestResult, error) {
	if es == nil {
		return nil, fmt.Errorf("digest: envset is nil")
	}

	targetKeys := keys
	if len(targetKeys) == 0 {
		for k := range es.Vars {
			targetKeys = append(targetKeys, k)
		}
	}
	sort.Strings(targetKeys)

	entries := make(map[string]string, len(targetKeys))
	h := sha256.New()

	for _, k := range targetKeys {
		v, ok := es.Vars[k]
		if !ok {
			return nil, fmt.Errorf("digest: key %q not found", k)
		}
		line := fmt.Sprintf("%s=%s\n", k, v)
		h.Write([]byte(line))
		kh := sha256.Sum256([]byte(line))
		entries[k] = hex.EncodeToString(kh[:])
	}

	return &DigestResult{
		Digest:  hex.EncodeToString(h.Sum(nil)),
		Entries: entries,
	}, nil
}

// DigestsMatch returns true when two EnvSets produce identical digests for the
// given key selection.
func DigestsMatch(a, b *EnvSet, keys []string) (bool, error) {
	da, err := Digest(a, keys)
	if err != nil {
		return false, fmt.Errorf("digest a: %w", err)
	}
	db, err := Digest(b, keys)
	if err != nil {
		return false, fmt.Errorf("digest b: %w", err)
	}
	return da.Digest == db.Digest, nil
}
