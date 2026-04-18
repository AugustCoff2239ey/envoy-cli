package envset

import (
	"errors"
	"fmt"
	"time"
)

// LockEntry represents a locked key in an EnvSet.
type LockEntry struct {
	Key       string    `json:"key"`
	LockedAt  time.Time `json:"locked_at"`
	LockedBy  string    `json:"locked_by"`
}

// ErrKeyLocked is returned when an operation attempts to modify a locked key.
var ErrKeyLocked = errors.New("key is locked and cannot be modified")

// LockKey marks a key in the EnvSet as locked, preventing modification.
// The lockedBy string identifies who or what locked the key.
func LockKey(es *EnvSet, key, lockedBy string) error {
	if es == nil {
		return errors.New("envset is nil")
	}
	if err := ValidateKey(key); err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}
	if _, exists := es.Vars[key]; !exists {
		return fmt.Errorf("key %q does not exist in envset", key)
	}
	if es.Locks == nil {
		es.Locks = make(map[string]LockEntry)
	}
	es.Locks[key] = LockEntry{
		Key:      key,
		LockedAt: time.Now().UTC(),
		LockedBy: lockedBy,
	}
	return nil
}

// UnlockKey removes the lock from a key in the EnvSet.
func UnlockKey(es *EnvSet, key string) error {
	if es == nil {
		return errors.New("envset is nil")
	}
	if es.Locks == nil {
		return fmt.Errorf("key %q is not locked", key)
	}
	if _, locked := es.Locks[key]; !locked {
		return fmt.Errorf("key %q is not locked", key)
	}
	delete(es.Locks, key)
	return nil
}

// IsLocked reports whether a key is currently locked.
func IsLocked(es *EnvSet, key string) bool {
	if es == nil || es.Locks == nil {
		return false
	}
	_, locked := es.Locks[key]
	return locked
}

// LockedKeys returns all currently locked keys in the EnvSet.
func LockedKeys(es *EnvSet) []LockEntry {
	if es == nil || es.Locks == nil {
		return nil
	}
	entries := make([]LockEntry, 0, len(es.Locks))
	for _, entry := range es.Locks {
		entries = append(entries, entry)
	}
	return entries
}

// GetLock returns the LockEntry for a key, or an error if the key is not locked.
func GetLock(es *EnvSet, key string) (LockEntry, error) {
	if es == nil {
		return LockEntry{}, errors.New("envset is nil")
	}
	if es.Locks == nil {
		return LockEntry{}, fmt.Errorf("key %q is not locked", key)
	}
	entry, locked := es.Locks[key]
	if !locked {
		return LockEntry{}, fmt.Errorf("key %q is not locked", key)
	}
	return entry, nil
}
