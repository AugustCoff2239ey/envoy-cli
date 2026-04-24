package envset

import (
	"fmt"
	"strings"
)

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	// MaxLength is the maximum allowed value length (default 64).
	MaxLength int
	// Suffix is appended to truncated values (default "...").
	Suffix string
	// Keys restricts truncation to specific keys; empty means all keys.
	Keys []string
}

// DefaultTruncateOptions returns sensible defaults for TruncateOptions.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength: 64,
		Suffix:    "...",
	}
}

// Truncate shortens values in es that exceed opts.MaxLength, appending opts.Suffix.
// If opts.Keys is non-empty only those keys are considered.
// Returns an error if es is nil or opts.MaxLength is not positive.
func Truncate(es *EnvSet, opts TruncateOptions) (map[string]string, error) {
	if es == nil {
		return nil, fmt.Errorf("truncate: nil EnvSet")
	}
	if opts.MaxLength <= 0 {
		return nil, fmt.Errorf("truncate: MaxLength must be positive, got %d", opts.MaxLength)
	}

	targetKeys := make(map[string]bool)
	for _, k := range opts.Keys {
		targetKeys[k] = true
	}

	result := make(map[string]string)

	for k, v := range es.Vars {
		if len(targetKeys) > 0 && !targetKeys[k] {
			continue
		}
		if len(v) > opts.MaxLength {
			cutoff := opts.MaxLength - len(opts.Suffix)
			if cutoff < 0 {
				cutoff = 0
			}
			result[k] = v[:cutoff] + opts.Suffix
		} else {
			result[k] = v
		}
	}

	return result, nil
}

// TruncateApply applies Truncate in-place, updating es.Vars with truncated values.
func TruncateApply(es *EnvSet, opts TruncateOptions) error {
	truncated, err := Truncate(es, opts)
	if err != nil {
		return err
	}
	for k, v := range truncated {
		if strings.Contains(v, opts.Suffix) || len(es.Vars[k]) <= opts.MaxLength {
			es.Vars[k] = v
		}
	}
	return nil
}
