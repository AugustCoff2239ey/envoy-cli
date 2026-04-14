package envset

import (
	"fmt"
	"time"
)

// PinEntry records a pinned value for a key at a specific version/timestamp.
type PinEntry struct {
	Key       string
	Value     string
	PinnedAt  time.Time
	PinnedBy  string
}

// PinKey pins the current value of a key in the EnvSet, preventing it from
// being overwritten by sync, merge, or promote operations unless explicitly unpinned.
func PinKey(es *EnvSet, key, pinnedBy string) (PinEntry, error) {
	if es == nil {
		return PinEntry{}, fmt.Errorf("pin: nil EnvSet")
	}
	if err := ValidateKey(key); err != nil {
		return PinEntry{}, fmt.Errorf("pin: %w", err)
	}
	val, ok := es.Vars[key]
	if !ok {
		return PinEntry{}, fmt.Errorf("pin: key %q not found", key)
	}
	if es.Meta == nil {
		es.Meta = make(map[string]string)
	}
	es.Meta["pinned:"+key] = "true"
	return PinEntry{
		Key:      key,
		Value:    val,
		PinnedAt: time.Now().UTC(),
		PinnedBy: pinnedBy,
	}, nil
}

// UnpinKey removes the pin from a key, allowing it to be modified freely.
func UnpinKey(es *EnvSet, key string) error {
	if es == nil {
		return fmt.Errorf("unpin: nil EnvSet")
	}
	if err := ValidateKey(key); err != nil {
		return fmt.Errorf("unpin: %w", err)
	}
	if es.Meta == nil || es.Meta["pinned:"+key] != "true" {
		return fmt.Errorf("unpin: key %q is not pinned", key)
	}
	delete(es.Meta, "pinned:"+key)
	return nil
}

// IsPinned reports whether the given key is currently pinned.
func IsPinned(es *EnvSet, key string) bool {
	if es == nil || es.Meta == nil {
		return false
	}
	return es.Meta["pinned:"+key] == "true"
}

// PinnedKeys returns all keys that are currently pinned in the EnvSet.
func PinnedKeys(es *EnvSet) []string {
	if es == nil || es.Meta == nil {
		return nil
	}
	var keys []string
	prefix := "pinned:"
	for k, v := range es.Meta {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix && v == "true" {
			keys = append(keys, k[len(prefix):])
		}
	}
	return keys
}
