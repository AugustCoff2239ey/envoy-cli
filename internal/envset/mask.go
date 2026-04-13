package envset

import (
	"errors"
	"regexp"
	"strings"
)

var sensitivePattern = regexp.MustCompile(
	`(?i)(password|secret|token|key|api|auth|credential|private|passphrase)`,
)

// MaskValue replaces all but the last 4 characters of a value with asterisks.
// If the value is 4 characters or shorter, the entire value is masked.
func MaskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return strings.Repeat("*", len(value)-4) + value[len(value)-4:]
}

// IsSensitiveKey returns true if the key name matches common sensitive patterns.
func IsSensitiveKey(key string) bool {
	return sensitivePattern.MatchString(key)
}

// MaskSensitive returns a copy of the EnvSet with sensitive values masked.
// Only keys matching sensitive naming patterns are masked.
func MaskSensitive(es *EnvSet) (*EnvSet, error) {
	if es == nil {
		return nil, errors.New("envset: cannot mask nil EnvSet")
	}

	masked := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
	}

	for k, v := range es.Vars {
		if IsSensitiveKey(k) {
			masked.Vars[k] = MaskValue(v)
		} else {
			masked.Vars[k] = v
		}
	}

	return masked, nil
}

// MaskKeys returns a copy of the EnvSet with the specified keys masked,
// regardless of whether they match sensitive patterns.
func MaskKeys(es *EnvSet, keys []string) (*EnvSet, error) {
	if es == nil {
		return nil, errors.New("envset: cannot mask nil EnvSet")
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	masked := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
	}

	for k, v := range es.Vars {
		if _, ok := keySet[k]; ok {
			masked.Vars[k] = MaskValue(v)
		} else {
			masked.Vars[k] = v
		}
	}

	return masked, nil
}
