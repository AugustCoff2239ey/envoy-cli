package envset

import (
	"fmt"
	"strings"
)

// RedactedEntry holds a key and its redacted representation.
type RedactedEntry struct {
	Key      string
	Redacted string
}

// RedactResult holds the redacted view of an EnvSet.
type RedactResult struct {
	Name        string
	Environment string
	Entries     []RedactedEntry
}

// Redact returns a RedactResult where all values are replaced with a
// placeholder, optionally revealing a configurable number of trailing chars.
// revealSuffix controls how many trailing characters of the value are shown.
// Pass 0 to hide the value entirely.
func Redact(es *EnvSet, revealSuffix int) (*RedactResult, error) {
	if es == nil {
		return nil, fmt.Errorf("redact: envset must not be nil")
	}
	if revealSuffix < 0 {
		return nil, fmt.Errorf("redact: revealSuffix must be >= 0, got %d", revealSuffix)
	}

	result := &RedactResult{
		Name:        es.Name,
		Environment: es.Environment,
	}

	for _, key := range sortedKeys(es.Vars) {
		val := es.Vars[key]
		result.Entries = append(result.Entries, RedactedEntry{
			Key:      key,
			Redacted: redactValue(val, revealSuffix),
		})
	}

	return result, nil
}

// redactValue replaces a value with asterisks, optionally revealing a suffix.
func redactValue(val string, revealSuffix int) string {
	if len(val) == 0 {
		return "[empty]"
	}
	if revealSuffix == 0 || revealSuffix >= len(val) {
		return strings.Repeat("*", len(val))
	}
	masked := strings.Repeat("*", len(val)-revealSuffix)
	return masked + val[len(val)-revealSuffix:]
}

// RedactKeys returns a new EnvSet with the values of the given keys replaced
// by a fixed placeholder string, leaving all other values intact.
func RedactKeys(es *EnvSet, keys []string, placeholder string) (*EnvSet, error) {
	if es == nil {
		return nil, fmt.Errorf("redact: envset must not be nil")
	}
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	redactSet := map[string]bool{}
	for _, k := range keys {
		redactSet[k] = true
	}

	out := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
	}
	for k, v := range es.Vars {
		if redactSet[k] {
			out.Vars[k] = placeholder
		} else {
			out.Vars[k] = v
		}
	}
	return out, nil
}
