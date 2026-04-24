package envset

import (
	"fmt"
	"strings"
)

// RevealResult holds the result of a reveal operation for a single key.
type RevealResult struct {
	Key     string
	Masked  string
	Visible string
}

// RevealOptions controls how values are partially revealed.
type RevealOptions struct {
	// SuffixLen is the number of trailing characters to show (default 4).
	SuffixLen int
	// PrefixLen is the number of leading characters to show (default 0).
	PrefixLen int
}

// DefaultRevealOptions returns sensible defaults for RevealOptions.
func DefaultRevealOptions() RevealOptions {
	return RevealOptions{SuffixLen: 4, PrefixLen: 0}
}

// Reveal partially exposes the value of specified keys according to RevealOptions.
// It returns a slice of RevealResult for each requested key.
func Reveal(es *EnvSet, keys []string, opts RevealOptions) ([]RevealResult, error) {
	if es == nil {
		return nil, ErrNilEnvSet
	}
	if opts.SuffixLen < 0 {
		return nil, fmt.Errorf("envset: SuffixLen must be >= 0")
	}
	if opts.PrefixLen < 0 {
		return nil, fmt.Errorf("envset: PrefixLen must be >= 0")
	}

	results := make([]RevealResult, 0, len(keys))
	for _, k := range keys {
		v, ok := es.Vars[k]
		if !ok {
			return nil, fmt.Errorf("envset: key %q not found", k)
		}
		visible := partialReveal(v, opts.PrefixLen, opts.SuffixLen)
		results = append(results, RevealResult{
			Key:     k,
			Masked:  MaskValue(v),
			Visible: visible,
		})
	}
	return results, nil
}

func partialReveal(value string, prefixLen, suffixLen int) string {
	n := len(value)
	shown := prefixLen + suffixLen
	if shown >= n {
		return value
	}
	prefix := value[:prefixLen]
	suffix := value[n-suffixLen:]
	midLen := n - shown
	return prefix + strings.Repeat("*", midLen) + suffix
}
