package envset

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyProtected = errors.New("key is already protected")
	ErrNotProtected     = errors.New("key is not protected")
)

const protectedMeta = "__protected__"

// ProtectKey marks a key as protected, preventing deletion or overwrite.
func ProtectKey(e *EnvSet, key string) error {
	if e == nil {
		return errors.New("nil EnvSet")
	}
	if err := ValidateKey(key); err != nil {
		return err
	}
	if _, ok := e.Vars[key]; !ok {
		return fmt.Errorf("key %q does not exist", key)
	}
	if e.Meta == nil {
		e.Meta = map[string]map[string]string{}
	}
	if e.Meta[key] == nil {
		e.Meta[key] = map[string]string{}
	}
	if e.Meta[key][protectedMeta] == "true" {
		return ErrAlreadyProtected
	}
	e.Meta[key][protectedMeta] = "true"
	return nil
}

// UnprotectKey removes protection from a key.
func UnprotectKey(e *EnvSet, key string) error {
	if e == nil {
		return errors.New("nil EnvSet")
	}
	if !IsProtected(e, key) {
		return ErrNotProtected
	}
	delete(e.Meta[key], protectedMeta)
	return nil
}

// IsProtected reports whether a key is protected.
func IsProtected(e *EnvSet, key string) bool {
	if e == nil || e.Meta == nil {
		return false
	}
	return e.Meta[key][protectedMeta] == "true"
}

// ProtectedKeys returns all protected keys in the EnvSet.
func ProtectedKeys(e *EnvSet) []string {
	var keys []string
	if e == nil || e.Meta == nil {
		return keys
	}
	for k := range e.Vars {
		if IsProtected(e, k) {
			keys = append(keys, k)
		}
	}
	return keys
}
