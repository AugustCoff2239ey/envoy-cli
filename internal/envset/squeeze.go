package envset

import "fmt"

// SqueezeOptions controls how duplicate values are handled.
type SqueezeOptions struct {
	// KeepFirst retains the first key encountered with a given value.
	// If false, the last key is kept.
	KeepFirst bool
}

// DefaultSqueezeOptions returns sensible defaults.
func DefaultSqueezeOptions() SqueezeOptions {
	return SqueezeOptions{KeepFirst: true}
}

// SqueezeResult holds keys removed during a squeeze operation.
type SqueezeResult struct {
	Removed []string
}

// Squeeze removes keys that share duplicate values, keeping one representative
// per unique value according to the provided options.
func Squeeze(es *EnvSet, opts SqueezeOptions) (*SqueezeResult, error) {
	if es == nil {
		return nil, fmt.Errorf("squeeze: nil EnvSet")
	}

	seen := make(map[string]string) // value -> key already kept
	var removed []string

	keys := sortedKeys(es.Vars)
	if !opts.KeepFirst {
		// reverse so last key wins when we iterate forward
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	for _, k := range keys {
		v := es.Vars[k]
		if existing, dup := seen[v]; dup {
			_ = existing
			removed = append(removed, k)
			delete(es.Vars, k)
		} else {
			seen[v] = k
		}
	}

	return &SqueezeResult{Removed: removed}, nil
}
