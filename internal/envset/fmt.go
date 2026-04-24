package envset

import (
	"fmt"
	"strings"
)

// FmtOptions controls how Format behaves.
type FmtOptions struct {
	// SortKeys sorts keys alphabetically before formatting.
	SortKeys bool
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// QuoteValues wraps all values in double quotes.
	QuoteValues bool
}

// DefaultFmtOptions returns sensible defaults for Format.
func DefaultFmtOptions() FmtOptions {
	return FmtOptions{
		SortKeys:   true,
		TrimValues: true,
	}
}

// Fmt applies formatting rules to the keys and values of an EnvSet in-place.
// It returns an error if es is nil or if any resulting key would be invalid.
func Fmt(es *EnvSet, opts FmtOptions) error {
	if es == nil {
		return fmt.Errorf("fmt: nil EnvSet")
	}
	if err := AssertMutable(es); err != nil {
		return fmt.Errorf("fmt: %w", err)
	}
	if err := AssertWritable(es); err != nil {
		return fmt.Errorf("fmt: %w", err)
	}

	keys := es.Keys()
	if opts.SortKeys {
		keys = SortedKeys(es)
	}

	// Collect transformed pairs first to avoid mid-iteration mutation issues.
	type pair struct{ k, v string }
	pairs := make([]pair, 0, len(keys))

	for _, k := range keys {
		v, _ := es.Get(k)

		newKey := k
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(k)
		}
		if err := ValidateKey(newKey); err != nil {
			return fmt.Errorf("fmt: invalid key after transform %q: %w", newKey, err)
		}

		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.QuoteValues {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		pairs = append(pairs, pair{newKey, v})
	}

	// Clear and repopulate.
	for _, k := range es.Keys() {
		es.Delete(k)
	}
	for _, p := range pairs {
		es.Set(p.k, p.v) //nolint:errcheck
	}
	return nil
}
