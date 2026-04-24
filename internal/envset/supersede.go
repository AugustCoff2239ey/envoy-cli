package envset

import "fmt"

// SupersedeOptions controls how Supersede behaves.
type SupersedeOptions struct {
	// Keys restricts superseding to specific keys; empty means all keys.
	Keys []string
	// SkipLocked prevents overwriting locked keys in the target.
	SkipLocked bool
	// SkipProtected prevents overwriting protected keys in the target.
	SkipProtected bool
}

// DefaultSupersedeOptions returns sensible defaults.
func DefaultSupersedeOptions() SupersedeOptions {
	return SupersedeOptions{
		SkipLocked:    true,
		SkipProtected: true,
	}
}

// Supersede overwrites keys in target with values from source, returning the
// number of keys that were actually replaced. Unlike Merge or Copy, Supersede
// always overwrites existing values (subject to lock/protect guards).
func Supersede(source, target *EnvSet, opts SupersedeOptions) (int, error) {
	if source == nil {
		return 0, fmt.Errorf("supersede: source is nil")
	}
	if target == nil {
		return 0, fmt.Errorf("supersede: target is nil")
	}
	if err := AssertMutable(target); err != nil {
		return 0, fmt.Errorf("supersede: %w", err)
	}
	if err := AssertWritable(target); err != nil {
		return 0, fmt.Errorf("supersede: %w", err)
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range source.Vars {
			keys = append(keys, k)
		}
	}

	replaced := 0
	for _, k := range keys {
		val, ok := source.Vars[k]
		if !ok {
			return replaced, fmt.Errorf("supersede: key %q not found in source", k)
		}
		if opts.SkipLocked && IsLocked(target, k) {
			continue
		}
		if opts.SkipProtected && IsProtected(target, k) {
			continue
		}
		target.Vars[k] = val
		replaced++
	}
	return replaced, nil
}
