package envset

import "fmt"

// MirrorOptions controls how a mirror operation behaves.
type MirrorOptions struct {
	// Keys to mirror; if empty, all keys are mirrored.
	Keys []string
	// Overwrite allows existing keys in the target to be replaced.
	Overwrite bool
	// Prefix is prepended to every mirrored key name in the target.
	Prefix string
}

// DefaultMirrorOptions returns sensible defaults for Mirror.
func DefaultMirrorOptions() MirrorOptions {
	return MirrorOptions{
		Overwrite: false,
	}
}

// Mirror copies key-value pairs from src into dst, optionally adding a prefix
// and respecting lock / protect guards on the destination set.
// It returns the number of keys successfully mirrored and any error encountered.
func Mirror(src, dst *EnvSet, opts MirrorOptions) (int, error) {
	if src == nil {
		return 0, fmt.Errorf("mirror: source EnvSet is nil")
	}
	if dst == nil {
		return 0, fmt.Errorf("mirror: destination EnvSet is nil")
	}
	if err := AssertMutable(dst); err != nil {
		return 0, fmt.Errorf("mirror: %w", err)
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src.Vars {
			keys = append(keys, k)
		}
	}

	count := 0
	for _, k := range keys {
		v, ok := src.Vars[k]
		if !ok {
			return count, fmt.Errorf("mirror: key %q not found in source", k)
		}

		destKey := opts.Prefix + k
		if err := ValidateKey(destKey); err != nil {
			return count, fmt.Errorf("mirror: invalid destination key %q: %w", destKey, err)
		}
		if IsLocked(dst, destKey) {
			continue
		}
		if IsProtected(dst, destKey) {
			continue
		}
		if _, exists := dst.Vars[destKey]; exists && !opts.Overwrite {
			continue
		}
		dst.Vars[destKey] = v
		count++
	}
	return count, nil
}
