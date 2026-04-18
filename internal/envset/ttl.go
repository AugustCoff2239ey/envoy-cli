package envset

import (
	"fmt"
	"time"
)

// SetTTL sets a time-to-live duration on a key, after which it is considered expired.
func SetTTL(es *EnvSet, key string, ttl time.Duration) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if _, ok := es.Vars[key]; !ok {
		return fmt.Errorf("key %q does not exist", key)
	}
	if err := ValidateKey(key); err != nil {
		return err
	}
	if ttl <= 0 {
		return fmt.Errorf("ttl must be positive")
	}
	expiry := time.Now().Add(ttl)
	return SetExpiry(es, key, expiry)
}

// GetTTL returns the remaining TTL for a key, or an error if not set or already expired.
func GetTTL(es *EnvSet, key string) (time.Duration, error) {
	if es == nil {
		return 0, fmt.Errorf("envset is nil")
	}
	expiry, err := GetExpiry(es, key)
	if err != nil {
		return 0, err
	}
	remaining := time.Until(expiry)
	if remaining <= 0 {
		return 0, fmt.Errorf("key %q has already expired", key)
	}
	return remaining, nil
}

// PurgeTTL removes all keys whose TTL has elapsed.
func PurgeTTL(es *EnvSet) ([]string, error) {
	return PurgeExpired(es)
}
