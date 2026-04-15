package envset

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenResult holds the output of a Flatten operation.
type FlattenResult struct {
	// Vars is the flattened map of key=value pairs.
	Vars map[string]string
	// Conflicts lists keys that had differing values across sources.
	Conflicts []string
}

// FlattenOptions controls how Flatten merges multiple EnvSets.
type FlattenOptions struct {
	// Prefix optionally prepends a string to every key.
	Prefix string
	// Overwrite determines whether later sets overwrite earlier ones on conflict.
	Overwrite bool
	// UppercaseKeys normalises all keys to uppercase.
	UppercaseKeys bool
}

// Flatten merges one or more EnvSets into a single key=value map.
// Sets are applied in order; conflict behaviour is governed by opts.
// Returns an error if any source EnvSet is nil.
func Flatten(opts FlattenOptions, sets ...*EnvSet) (*FlattenResult, error) {
	if len(sets) == 0 {
		return &FlattenResult{Vars: map[string]string{}}, nil
	}

	result := &FlattenResult{
		Vars: make(map[string]string),
	}

	conflictSet := make(map[string]bool)

	for i, es := range sets {
		if es == nil {
			return nil, fmt.Errorf("flatten: source at index %d is nil", i)
		}

		for k, v := range es.Vars {
			key := k
			if opts.UppercaseKeys {
				key = strings.ToUpper(key)
			}
			if opts.Prefix != "" {
				key = opts.Prefix + key
			}

			if existing, exists := result.Vars[key]; exists {
				if existing != v {
					conflictSet[key] = true
				}
				if !opts.Overwrite {
					continue
				}
			}
			result.Vars[key] = v
		}
	}

	for k := range conflictSet {
		result.Conflicts = append(result.Conflicts, k)
	}
	sort.Strings(result.Conflicts)

	return result, nil
}
