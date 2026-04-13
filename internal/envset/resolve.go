package envset

import "fmt"

// ResolveResult holds the outcome of resolving a single variable reference.
type ResolveResult struct {
	Key      string
	Resolved string
	Source   string // "local", "override", or "default"
}

// ResolveOptions controls how variable resolution behaves.
type ResolveOptions struct {
	// Overrides are key-value pairs that take precedence over the EnvSet values.
	Overrides map[string]string
	// Defaults are key-value pairs used when a key is absent from both the
	// EnvSet and Overrides.
	Defaults map[string]string
}

// Resolve returns the effective value for every key in the EnvSet, applying
// overrides and falling back to defaults where necessary.
//
// Resolution order (highest to lowest priority):
//  1. opts.Overrides
//  2. EnvSet vars
//  3. opts.Defaults
func Resolve(es *EnvSet, opts ResolveOptions) ([]ResolveResult, error) {
	if es == nil {
		return nil, fmt.Errorf("resolve: envset must not be nil")
	}

	// Collect the union of all keys.
	keySet := make(map[string]struct{})
	for k := range es.Vars {
		keySet[k] = struct{}{}
	}
	for k := range opts.Overrides {
		keySet[k] = struct{}{}
	}
	for k := range opts.Defaults {
		keySet[k] = struct{}{}
	}

	results := make([]ResolveResult, 0, len(keySet))
	for k := range keySet {
		if err := ValidateKey(k); err != nil {
			return nil, fmt.Errorf("resolve: invalid key %q: %w", k, err)
		}

		var resolved, source string
		if v, ok := opts.Overrides[k]; ok {
			resolved = v
			source = "override"
		} else if v, ok := es.Vars[k]; ok {
			resolved = v
			source = "local"
		} else if v, ok := opts.Defaults[k]; ok {
			resolved = v
			source = "default"
		}

		results = append(results, ResolveResult{
			Key:      k,
			Resolved: resolved,
			Source:   source,
		})
	}

	return results, nil
}

// ResolveKey resolves the effective value for a single key.
func ResolveKey(es *EnvSet, key string, opts ResolveOptions) (ResolveResult, error) {
	if es == nil {
		return ResolveResult{}, fmt.Errorf("resolve: envset must not be nil")
	}
	if err := ValidateKey(key); err != nil {
		return ResolveResult{}, fmt.Errorf("resolve: %w", err)
	}

	if v, ok := opts.Overrides[key]; ok {
		return ResolveResult{Key: key, Resolved: v, Source: "override"}, nil
	}
	if v, ok := es.Vars[key]; ok {
		return ResolveResult{Key: key, Resolved: v, Source: "local"}, nil
	}
	if v, ok := opts.Defaults[key]; ok {
		return ResolveResult{Key: key, Resolved: v, Source: "default"}, nil
	}

	return ResolveResult{}, fmt.Errorf("resolve: key %q not found", key)
}
