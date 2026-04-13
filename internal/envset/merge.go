package envset

import "fmt"

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the base value on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the incoming value on conflict.
	MergeStrategyTheirs
	// MergeStrategyError returns an error on conflict.
	MergeStrategyError
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged   *EnvSet
	Conflicts []string
}

// Merge combines the keys from src into dst according to the given strategy.
// dst is not mutated; a new EnvSet is returned inside MergeResult.
func Merge(dst, src *EnvSet, strategy MergeStrategy) (*MergeResult, error) {
	if dst == nil || src == nil {
		return nil, fmt.Errorf("merge: dst and src must not be nil")
	}

	merged, err := New(dst.Name, dst.Environment)
	if err != nil {
		return nil, fmt.Errorf("merge: %w", err)
	}

	// Copy all keys from dst first.
	for k, v := range dst.vars {
		_ = merged.Set(k, v)
	}

	var conflicts []string

	for k, srcVal := range src.vars {
		dstVal, exists := merged.vars[k]
		if !exists {
			_ = merged.Set(k, srcVal)
			continue
		}
		if dstVal == srcVal {
			continue
		}
		// Conflict detected.
		conflicts = append(conflicts, k)
		switch strategy {
		case MergeStrategyOurs:
			// keep dstVal — already set
		case MergeStrategyTheirs:
			_ = merged.Set(k, srcVal)
		case MergeStrategyError:
			return nil, fmt.Errorf("merge conflict on key %q: %q vs %q", k, dstVal, srcVal)
		}
	}

	return &MergeResult{Merged: merged, Conflicts: conflicts}, nil
}
