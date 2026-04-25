package envset

import "fmt"

// ShrinkOptions controls which keys are removed during shrinking.
type ShrinkOptions struct {
	// MaxKeys is the maximum number of keys to retain. 0 means no limit.
	MaxKeys int
	// KeepKeys is an explicit list of keys to always retain.
	KeepKeys []string
	// RemoveExpired removes keys that have passed their expiry time.
	RemoveExpired bool
	// RemoveEmpty removes keys whose values are empty strings.
	RemoveEmpty bool
}

// DefaultShrinkOptions returns sensible defaults for Shrink.
func DefaultShrinkOptions() ShrinkOptions {
	return ShrinkOptions{
		MaxKeys:       0,
		KeepKeys:      nil,
		RemoveExpired: true,
		RemoveEmpty:   false,
	}
}

// Shrink reduces the size of an EnvSet by removing keys according to opts.
// Keys listed in opts.KeepKeys are never removed. If MaxKeys > 0, keys are
// dropped (in insertion order) until the set is within the limit.
func Shrink(es *EnvSet, opts ShrinkOptions) error {
	if es == nil {
		return fmt.Errorf("shrink: envset is nil")
	}

	keep := make(map[string]bool, len(opts.KeepKeys))
	for _, k := range opts.KeepKeys {
		keep[k] = true
	}

	// Remove expired keys first.
	if opts.RemoveExpired {
		for _, k := range es.Keys() {
			if keep[k] {
				continue
			}
			if IsExpired(es, k) {
				if err := es.Delete(k); err != nil {
					return fmt.Errorf("shrink: removing expired key %q: %w", k, err)
				}
			}
		}
	}

	// Remove empty-value keys.
	if opts.RemoveEmpty {
		for _, k := range es.Keys() {
			if keep[k] {
				continue
			}
			if v, _ := es.Get(k); v == "" {
				if err := es.Delete(k); err != nil {
					return fmt.Errorf("shrink: removing empty key %q: %w", k, err)
				}
			}
		}
	}

	// Enforce MaxKeys limit.
	if opts.MaxKeys > 0 {
		keys := es.Keys()
		for len(keys) > opts.MaxKeys {
			candidate := ""
			for _, k := range keys {
				if !keep[k] {
					candidate = k
					break
				}
			}
			if candidate == "" {
				return fmt.Errorf("shrink: cannot reduce below %d keys; all remaining keys are protected by KeepKeys", len(keys))
			}
			if err := es.Delete(candidate); err != nil {
				return fmt.Errorf("shrink: removing key %q: %w", candidate, err)
			}
			keys = es.Keys()
		}
	}

	return nil
}
