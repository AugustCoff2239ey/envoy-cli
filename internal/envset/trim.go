package envset

import (
	"fmt"
	"strings"
)

// TrimResult holds the outcome of a trim operation.
type TrimResult struct {
	Key     string
	OldVal  string
	NewVal  string
	Changed bool
}

// TrimOptions controls which transformations Trim applies.
type TrimOptions struct {
	LeadingSpace  bool
	TrailingSpace bool
	Quotes        bool
	Keys          []string // if empty, all keys are processed
}

// DefaultTrimOptions returns TrimOptions with sensible defaults.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		LeadingSpace:  true,
		TrailingSpace: true,
		Quotes:        false,
	}
}

// Trim cleans up whitespace (and optionally surrounding quotes) from env var
// values in es. It returns a slice of TrimResult describing every changed key.
func Trim(es *EnvSet, opts TrimOptions) ([]TrimResult, error) {
	if es == nil {
		return nil, fmt.Errorf("trim: envset is nil")
	}

	targets := opts.Keys
	if len(targets) == 0 {
		for k := range es.Vars {
			targets = append(targets, k)
		}
	}

	var results []TrimResult

	for _, k := range targets {
		v, ok := es.Vars[k]
		if !ok {
			return nil, fmt.Errorf("trim: key %q not found", k)
		}

		newVal := v
		if opts.LeadingSpace {
			newVal = strings.TrimLeft(newVal, " \t")
		}
		if opts.TrailingSpace {
			newVal = strings.TrimRight(newVal, " \t")
		}
		if opts.Quotes {
			if len(newVal) >= 2 {
				if (newVal[0] == '"' && newVal[len(newVal)-1] == '"') ||
					(newVal[0] == '\'' && newVal[len(newVal)-1] == '\'') {
					newVal = newVal[1 : len(newVal)-1]
				}
			}
		}

		results = append(results, TrimResult{
			Key:     k,
			OldVal:  v,
			NewVal:  newVal,
			Changed: v != newVal,
		})

		if v != newVal {
			es.Vars[k] = newVal
		}
	}

	return results, nil
}
