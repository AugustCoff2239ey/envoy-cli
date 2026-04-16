package envset

import (
	"fmt"
	"time"
)

const (
	metaExpireAt = "__expire_at__"
)

// SetExpiry sets an expiration time on a key in the EnvSet.
func SetExpiry(e *EnvSet, key string, expiresAt time.Time) error {
	if e == nil {
		return fmt.Errorf("envset is nil")
	}
	if err := ValidateKey(key); err != nil {
		return err
	}
	if _, ok := e.Vars[key]; !ok {
		return fmt.Errorf("key %q does not exist", key)
	}
	if e.Meta == nil {
		e.Meta = map[string]string{}
	}
	e.Meta[metaExpireAt+"."+key] = expiresAt.UTC().Format(time.RFC3339)
	return nil
}

// GetExpiry returns the expiration time for a key, and whether one is set.
func GetExpiry(e *EnvSet, key string) (time.Time, bool, error) {
	if e == nil {
		return time.Time{}, false, fmt.Errorf("envset is nil")
	}
	raw, ok := e.Meta[metaExpireAt+"."+key]
	if !ok {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid expiry format for key %q: %w", key, err)
	}
	return t, true, nil
}

// IsExpired reports whether a key has passed its expiration time.
func IsExpired(e *EnvSet, key string) (bool, error) {
	t, set, err := GetExpiry(e, key)
	if err != nil || !set {
		return false, err
	}
	return time.Now().UTC().After(t), nil
}

// PurgeExpired removes all keys from the EnvSet that have expired.
func PurgeExpired(e *EnvSet) ([]string, error) {
	if e == nil {
		return nil, fmt.Errorf("envset is nil")
	}
	var purged []string
	for key := range e.Vars {
		expired, err := IsExpired(e, key)
		if err != nil {
			return purged, err
		}
		if expired {
			delete(e.Vars, key)
			delete(e.Meta, metaExpireAt+"."+key)
			purged = append(purged, key)
		}
	}
	return purged, nil
}
