package envset

import "fmt"

// DedupResult holds the outcome of a deduplication operation.
type DedupResult struct {
	Removed []string
	Kept    map[string]string
}

// DedupStrategy controls which value to keep when duplicates are detected
// across a list of EnvSets sharing the same key.
type DedupStrategy string

const (
	DedupKeepFirst DedupStrategy = "first"
	DedupKeepLast  DedupStrategy = "last"
)

// Dedup scans the given EnvSet for keys that appear in any of the provided
// reference sets and removes them, keeping only the canonical value according
// to the chosen strategy. When strategy is DedupKeepLast the value from the
// last reference set wins; DedupKeepFirst retains the value already present
// in base.
func Dedup(base *EnvSet, refs []*EnvSet, strategy DedupStrategy) (*DedupResult, error) {
	if base == nil {
		return nil, fmt.Errorf("dedup: base EnvSet must not be nil")
	}
	if strategy != DedupKeepFirst && strategy != DedupKeepLast {
		return nil, fmt.Errorf("dedup: unknown strategy %q", strategy)
	}

	result := &DedupResult{
		Kept: make(map[string]string),
	}

	for key, val := range base.Vars {
		result.Kept[key] = val
	}

	for _, ref := range refs {
		if ref == nil {
			continue
		}
		for key, refVal := range ref.Vars {
			if _, exists := result.Kept[key]; exists {
				if strategy == DedupKeepLast {
					result.Kept[key] = refVal
				}
				result.Removed = append(result.Removed, key)
			}
		}
	}

	// Apply kept values back to base
	for key, val := range result.Kept {
		base.Vars[key] = val
	}

	return result, nil
}
