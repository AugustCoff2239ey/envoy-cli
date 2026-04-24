package envset

import "fmt"

// CascadeOptions controls how values flow through the cascade chain.
type CascadeOptions struct {
	// Overwrite allows downstream sets to overwrite upstream values.
	// When false, upstream values take priority (first-wins).
	Overwrite bool

	// SkipLocked prevents locked keys from being overwritten even if Overwrite is true.
	SkipLocked bool

	// SkipProtected prevents protected keys from being overwritten even if Overwrite is true.
	SkipProtected bool
}

// DefaultCascadeOptions returns sensible defaults for cascading:
// upstream values win, locked and protected keys are always respected.
func DefaultCascadeOptions() CascadeOptions {
	return CascadeOptions{
		Overwrite:     false,
		SkipLocked:    true,
		SkipProtected: true,
	}
}

// CascadeResult holds the outcome of a Cascade operation.
type CascadeResult struct {
	// Applied is the number of keys written into the destination.
	Applied int
	// Skipped is the number of keys that were not written due to conflicts or guards.
	Skipped int
	// Sources records, per key, which source index provided the final value.
	Sources map[string]int
}

// Cascade merges values from one or more source EnvSets into dst following a
// layered priority model. Sources are evaluated left-to-right; the first source
// that defines a key wins unless opts.Overwrite is true, in which case later
// sources overwrite earlier ones (last-wins). Locked and protected keys in dst
// are guarded according to opts.
//
// This is useful for building environment configs from a hierarchy such as:
//
//	Cascade(dst, opts, defaults, shared, team, personal)
//
// Returns an error if dst or any source is nil.
func Cascade(dst *EnvSet, opts CascadeOptions, sources ...*EnvSet) (CascadeResult, error) {
	if dst == nil {
		return CascadeResult{}, fmt.Errorf("cascade: destination EnvSet must not be nil")
	}
	for i, src := range sources {
		if src == nil {
			return CascadeResult{}, fmt.Errorf("cascade: source at index %d is nil", i)
		}
	}

	result := CascadeResult{
		Sources: make(map[string]int),
	}

	for srcIdx, src := range sources {
		for key, val := range src.Vars {
			// Determine whether we should write this key into dst.
			_, exists := dst.Vars[key]

			if exists && !opts.Overwrite {
				// First-wins: upstream already has the key.
				result.Skipped++
				continue
			}

			if exists && opts.SkipLocked && IsLocked(dst, key) {
				result.Skipped++
				continue
			}

			if exists && opts.SkipProtected && IsProtected(dst, key) {
				result.Skipped++
				continue
			}

			dst.Vars[key] = val
			result.Sources[key] = srcIdx
			result.Applied++
		}
	}

	return result, nil
}
