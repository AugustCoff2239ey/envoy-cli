package envset

import (
	"strings"
)

var sensitiveKeyPatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE_KEY", "AUTH", "CREDENTIAL", "ACCESS_KEY",
}

// MaskValue replaces all but the last 4 characters of a value with asterisks.
// If the value is 4 characters or fewer, it is fully masked.
func MaskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return strings.Repeat("*", len(value)-4) + value[len(value)-4:]
}

// IsSensitiveKey returns true if the key matches common sensitive naming patterns.
func IsSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitiveKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskSensitive masks values for all keys that match sensitive naming patterns.
// Returns a new EnvSet with masked values; the original is unchanged.
func MaskSensitive(es *EnvSet) (*EnvSet, error) {
	if es == nil {
		return nil, ErrNilEnvSet
	}
	out := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
		Meta:        make(map[string]string, len(es.Meta)),
	}
	for k, v := range es.Vars {
		if IsSensitiveKey(k) {
			out.Vars[k] = MaskValue(v)
		} else {
			out.Vars[k] = v
		}
	}
	for k, v := range es.Meta {
		out.Meta[k] = v
	}
	return out, nil
}

// MaskKeys masks values for the specified keys only.
// Returns a new EnvSet with masked values; the original is unchanged.
func MaskKeys(es *EnvSet, keys []string) (*EnvSet, error) {
	if es == nil {
		return nil, ErrNilEnvSet
	}
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}
	out := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
		Meta:        make(map[string]string, len(es.Meta)),
	}
	for k, v := range es.Vars {
		if keySet[k] {
			out.Vars[k] = MaskValue(v)
		} else {
			out.Vars[k] = v
		}
	}
	for k, v := range es.Meta {
		out.Meta[k] = v
	}
	return out, nil
}
