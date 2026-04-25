package envset

import (
	"fmt"
	"sort"
)

// GradientEntry represents a single step in an environment gradient.
type GradientEntry struct {
	Environment string
	Key         string
	Value       string
}

// GradientResult holds the progression of a key's value across environments.
type GradientResult struct {
	Key     string
	Steps   []GradientEntry
	Uniform bool
}

// Gradient traces how one or more keys change across multiple EnvSet layers,
// ordered by their environment name. It is useful for auditing value drift
// across local → staging → production pipelines.
func Gradient(sets []*EnvSet, keys []string) ([]GradientResult, error) {
	if len(sets) == 0 {
		return nil, fmt.Errorf("gradient: at least one envset is required")
	}

	// Collect all keys if none specified.
	if len(keys) == 0 {
		keySet := map[string]struct{}{}
		for _, s := range sets {
			if s == nil {
				continue
			}
			for k := range s.Vars {
				keySet[k] = struct{}{}
			}
		}
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	results := make([]GradientResult, 0, len(keys))
	for _, key := range keys {
		if err := ValidateKey(key); err != nil {
			return nil, fmt.Errorf("gradient: invalid key %q: %w", key, err)
		}
		var steps []GradientEntry
		for _, s := range sets {
			if s == nil {
				continue
			}
			val, ok := s.Vars[key]
			if !ok {
				val = ""
			}
			steps = append(steps, GradientEntry{
				Environment: s.Environment,
				Key:         key,
				Value:       val,
			})
		}
		uniform := true
		for i := 1; i < len(steps); i++ {
			if steps[i].Value != steps[0].Value {
				uniform = false
				break
			}
		}
		results = append(results, GradientResult{
			Key:     key,
			Steps:   steps,
			Uniform: uniform,
		})
	}
	return results, nil
}
