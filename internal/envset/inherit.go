package envset

import "fmt"

// InheritResult holds the outcome of an inherit operation.
type InheritResult struct {
	Inherited []string
	Skipped   []string
}

// Inherit copies keys from parent into child for any key that does not already
// exist in child. Keys present in child are left untouched (skipped).
// If keys is non-empty only those keys are considered; otherwise all parent
// keys are candidates.
func Inherit(parent, child *EnvSet, keys []string) (*InheritResult, error) {
	if parent == nil || child == nil {
		return nil, fmt.Errorf("inherit: parent and child must not be nil")
	}

	candidates := keys
	if len(candidates) == 0 {
		for k := range parent.Vars {
			candidates = append(candidates, k)
		}
	}

	result := &InheritResult{}

	for _, k := range candidates {
		if err := ValidateKey(k); err != nil {
			return nil, fmt.Errorf("inherit: invalid key %q: %w", k, err)
		}

		parentVal, ok := parent.Vars[k]
		if !ok {
			return nil, fmt.Errorf("inherit: key %q not found in parent", k)
		}

		if _, exists := child.Vars[k]; exists {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		child.Vars[k] = parentVal
		result.Inherited = append(result.Inherited, k)
	}

	return result, nil
}
