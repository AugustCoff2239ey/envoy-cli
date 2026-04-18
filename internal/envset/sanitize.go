package envset

import (
	"fmt"
	"strings"
)

// SanitizeOptions controls sanitize behaviour.
type SanitizeOptions struct {
	StripControlChars bool
	TrimWhitespace    bool
	ReplaceSpacesInKeys bool
	SpaceReplacement    string
}

// DefaultSanitizeOptions returns sensible defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		StripControlChars:   true,
		TrimWhitespace:      true,
		ReplaceSpacesInKeys: true,
		SpaceReplacement:    "_",
	}
}

// SanitizeResult holds per-key changes made during sanitization.
type SanitizeResult struct {
	Key     string
	OldKey  string
	OldVal  string
	NewVal  string
	Changed bool
}

// Sanitize cleans keys and values of e according to opts.
// It returns a slice of SanitizeResult describing what changed.
func Sanitize(e *EnvSet, opts SanitizeOptions) ([]SanitizeResult, error) {
	if e == nil {
		return nil, fmt.Errorf("sanitize: nil EnvSet")
	}

	var results []SanitizeResult

	for k, v := range e.Vars {
		newKey := k
		newVal := v

		if opts.ReplaceSpacesInKeys {
			newKey = strings.ReplaceAll(newKey, " ", opts.SpaceReplacement)
		}

		if opts.TrimWhitespace {
			newVal = strings.TrimSpace(newVal)
		}

		if opts.StripControlChars {
			newVal = stripControl(newVal)
		}

		changed := newKey != k || newVal != v
		if changed {
			if newKey != k {
				delete(e.Vars, k)
			}
			e.Vars[newKey] = newVal
		}

		results = append(results, SanitizeResult{
			Key:     newKey,
			OldKey:  k,
			OldVal:  v,
			NewVal:  newVal,
			Changed: changed,
		})
	}

	return results, nil
}

func stripControl(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 32 || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
