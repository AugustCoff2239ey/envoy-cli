package envset

import (
	"fmt"
	"regexp"
)

// ExtractResult holds the result of an Extract operation.
type ExtractResult struct {
	Extracted *EnvSet
	Keys      []string
}

// ExtractOptions controls how keys are extracted.
type ExtractOptions struct {
	// Pattern is an optional regex to match keys.
	Pattern string
	// Keys is an explicit list of keys to extract.
	Keys []string
	// RemoveFromSource removes extracted keys from the source set.
	RemoveFromSource bool
}

// Extract copies matching keys from src into a new EnvSet.
// If RemoveFromSource is true, keys are deleted from src after extraction.
func Extract(src *EnvSet, opts ExtractOptions) (*ExtractResult, error) {
	if src == nil {
		return nil, fmt.Errorf("extract: source EnvSet is nil")
	}

	dest, err := New(src.Name, src.Environment)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	var matched []string

	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("extract: invalid pattern: %w", err)
		}
		for k, v := range src.Vars {
			if re.MatchString(k) {
				matched = append(matched, k)
				_ = v
			}
		}
	} else {
		for _, k := range opts.Keys {
			if _, ok := src.Vars[k]; !ok {
				return nil, fmt.Errorf("extract: key %q not found", k)
			}
			matched = append(matched, k)
		}
	}

	for _, k := range matched {
		if err := dest.Set(k, src.Vars[k]); err != nil {
			return nil, fmt.Errorf("extract: set key %q: %w", k, err)
		}
		if opts.RemoveFromSource {
			delete(src.Vars, k)
		}
	}

	return &ExtractResult{Extracted: dest, Keys: matched}, nil
}
