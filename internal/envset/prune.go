package envset

import "fmt"

// PruneOptions controls which keys are removed during pruning.
type PruneOptions struct {
	// RemoveEmpty removes keys with empty string values.
	RemoveEmpty bool
	// RemoveDuplicateValues removes keys whose values are duplicates of an
	// earlier key (first occurrence is kept).
	RemoveDuplicateValues bool
	// Keys is an explicit list of keys to remove. When non-empty, only these
	// keys are removed regardless of other options.
	Keys []string
}

// DefaultPruneOptions returns PruneOptions with sensible defaults.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		RemoveEmpty:           true,
		RemoveDuplicateValues: false,
	}
}

// Prune removes keys from es according to opts. It returns the number of keys
// removed and an error if es is nil or a key is locked/protected.
func Prune(es *EnvSet, opts PruneOptions) (int, error) {
	if es == nil {
		return 0, fmt.Errorf("prune: envset is nil")
	}

	// Explicit key list takes precedence.
	if len(opts.Keys) > 0 {
		removed := 0
		for _, k := range opts.Keys {
			if _, exists := es.Vars[k]; !exists {
				continue
			}
			if IsLocked(es, k) {
				return removed, fmt.Errorf("prune: key %q is locked", k)
			}
			if IsProtected(es, k) {
				return removed, fmt.Errorf("prune: key %q is protected", k)
			}
			delete(es.Vars, k)
			removed++
		}
		return removed, nil
	}

	seen := make(map[string]struct{})
	removed := 0

	for k, v := range es.Vars {
		if IsLocked(es, k) || IsProtected(es, k) {
			if opts.RemoveEmpty && v == "" {
				return removed, fmt.Errorf("prune: key %q is locked/protected and cannot be removed", k)
			}
			seen[v] = struct{}{}
			continue
		}

		if opts.RemoveEmpty && v == "" {
			delete(es.Vars, k)
			removed++
			continue
		}

		if opts.RemoveDuplicateValues {
			if _, dup := seen[v]; dup {
				delete(es.Vars, k)
				removed++
				continue
			}
			seen[v] = struct{}{}
		}
	}

	return removed, nil
}
