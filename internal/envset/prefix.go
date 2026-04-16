package envset

import (
	"fmt"
	"strings"
)

// AddPrefix adds a string prefix to all (or selected) keys in the EnvSet.
// If keys is empty, all keys are prefixed. Locked or protected keys are skipped.
func AddPrefix(es *EnvSet, prefix string, keys []string) (int, error) {
	if es == nil {
		return 0, fmt.Errorf("envset: nil EnvSet")
	}
	if prefix == "" {
		return 0, fmt.Errorf("envset: prefix must not be empty")
	}
	if err := ValidateKey(prefix + "A"); err != nil {
		return 0, fmt.Errorf("envset: invalid prefix %q: %w", prefix, err)
	}

	targets := keys
	if len(targets) == 0 {
		for k := range es.Vars {
			targets = append(targets, k)
		}
	}

	count := 0
	for _, k := range targets {
		val, ok := es.Vars[k]
		if !ok {
			return count, fmt.Errorf("envset: key %q not found", k)
		}
		if IsLocked(es, k) || IsProtected(es, k) {
			continue
		}
		newKey := prefix + k
		delete(es.Vars, k)
		es.Vars[newKey] = val
		count++
	}
	return count, nil
}

// StripPrefix removes a prefix from all (or selected) keys. Keys that do not
// carry the prefix are silently skipped.
func StripPrefix(es *EnvSet, prefix string, keys []string) (int, error) {
	if es == nil {
		return 0, fmt.Errorf("envset: nil EnvSet")
	}
	if prefix == "" {
		return 0, fmt.Errorf("envset: prefix must not be empty")
	}

	targets := keys
	if len(targets) == 0 {
		for k := range es.Vars {
			targets = append(targets, k)
		}
	}

	count := 0
	for _, k := range targets {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		val, ok := es.Vars[k]
		if !ok {
			continue
		}
		if IsLocked(es, k) || IsProtected(es, k) {
			continue
		}
		newKey := strings.TrimPrefix(k, prefix)
		if newKey == "" {
			continue
		}
		delete(es.Vars, k)
		es.Vars[newKey] = val
		count++
	}
	return count, nil
}
