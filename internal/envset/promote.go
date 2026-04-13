package envset

import "fmt"

// PromoteOptions configures how a promotion is performed.
type PromoteOptions struct {
	// Overwrite controls whether existing keys in the target are overwritten.
	Overwrite bool
	// Keys restricts promotion to a specific subset of keys. If empty, all keys are promoted.
	Keys []string
}

// Promote copies variables from a source EnvSet into a target EnvSet,
// optionally filtering by key and controlling overwrite behaviour.
// The source EnvSet is not modified. A new EnvSet representing the
// updated target is returned.
func Promote(source, target *EnvSet, opts PromoteOptions) (*EnvSet, error) {
	if source == nil {
		return nil, fmt.Errorf("promote: source EnvSet must not be nil")
	}
	if target == nil {
		return nil, fmt.Errorf("promote: target EnvSet must not be nil")
	}

	// Build the set of keys to promote.
	keys := opts.Keys
	if len(keys) == 0 {
		keys = make([]string, 0, len(source.Vars))
		for k := range source.Vars {
			keys = append(keys, k)
		}
	}

	// Validate that every requested key exists in the source.
	for _, k := range keys {
		if _, ok := source.Vars[k]; !ok {
			return nil, fmt.Errorf("promote: key %q not found in source %s/%s", k, source.Name, source.Environment)
		}
	}

	// Clone the target so we do not mutate the caller's value.
	result, err := Clone(target, "", "")
	if err != nil {
		return nil, fmt.Errorf("promote: %w", err)
	}

	for _, k := range keys {
		if _, exists := result.Vars[k]; exists && !opts.Overwrite {
			continue
		}
		result.Vars[k] = source.Vars[k]
	}

	return result, nil
}
