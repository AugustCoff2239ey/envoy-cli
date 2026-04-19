package envset

import "fmt"

// SplitResult holds the two EnvSets produced by a Split operation.
type SplitResult struct {
	Matched   *EnvSet
	Unmatched *EnvSet
}

// SplitOptions controls how Split partitions an EnvSet.
type SplitOptions struct {
	// Keys explicitly selected for the matched set. If empty, Predicate is used.
	Keys []string
	// Predicate is called for each key when Keys is empty.
	Predicate func(key, value string) bool
	// MatchedName overrides the name of the matched EnvSet.
	MatchedName string
	// UnmatchedName overrides the name of the unmatched EnvSet.
	UnmatchedName string
}

// Split partitions src into two EnvSets: one containing matched keys and one
// containing the remaining keys. The original EnvSet is not modified.
func Split(src *EnvSet, opts SplitOptions) (*SplitResult, error) {
	if src == nil {
		return nil, fmt.Errorf("split: source EnvSet is nil")
	}

	matchedName := opts.MatchedName
	if matchedName == "" {
		matchedName = src.Name + "-matched"
	}
	unmatchedName := opts.UnmatchedName
	if unmatchedName == "" {
		unmatchedName = src.Name + "-unmatched"
	}

	matched, err := New(matchedName, src.Environment)
	if err != nil {
		return nil, fmt.Errorf("split: %w", err)
	}
	unmatched, err := New(unmatchedName, src.Environment)
	if err != nil {
		return nil, fmt.Errorf("split: %w", err)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	for k, v := range src.Vars {
		var inMatch bool
		if len(opts.Keys) > 0 {
			_, inMatch = keySet[k]
		} else if opts.Predicate != nil {
			inMatch = opts.Predicate(k, v)
		}
		if inMatch {
			_ = matched.Set(k, v)
		} else {
			_ = unmatched.Set(k, v)
		}
	}

	return &SplitResult{Matched: matched, Unmatched: unmatched}, nil
}
