package envset

import (
	"fmt"
	"sort"
)

// ClipOptions controls how Clip trims an EnvSet down to a maximum number of keys.
type ClipOptions struct {
	// MaxKeys is the maximum number of keys to retain.
	MaxKeys int
	// Keys is an explicit ordered list of keys to keep; if empty, keys are sorted alphabetically.
	Keys []string
	// SkipLocked prevents locked keys from being removed even if they exceed MaxKeys.
	SkipLocked bool
}

// DefaultClipOptions returns sensible defaults for Clip.
func DefaultClipOptions() ClipOptions {
	return ClipOptions{
		MaxKeys:    10,
		SkipLocked: true,
	}
}

// Clip reduces an EnvSet to at most MaxKeys entries.
// If Keys is provided those are kept first (in order); remaining slots are
// filled alphabetically. Locked keys are always retained when SkipLocked is set.
func Clip(es *EnvSet, opts ClipOptions) ([]string, error) {
	if es == nil {
		return nil, fmt.Errorf("clip: nil EnvSet")
	}
	if opts.MaxKeys <= 0 {
		return nil, fmt.Errorf("clip: MaxKeys must be greater than zero")
	}

	kept := make(map[string]bool)
	var removed []string

	// Always keep explicitly requested keys first.
	for _, k := range opts.Keys {
		if _, ok := es.Vars[k]; ok {
			kept[k] = true
		}
	}

	// Fill remaining slots alphabetically.
	if len(kept) < opts.MaxKeys {
		all := make([]string, 0, len(es.Vars))
		for k := range es.Vars {
			if !kept[k] {
				all = append(all, k)
			}
		}
		sort.Strings(all)
		for _, k := range all {
			if len(kept) >= opts.MaxKeys {
				break
			}
			kept[k] = true
		}
	}

	// Remove keys not in the kept set, respecting SkipLocked.
	for k := range es.Vars {
		if kept[k] {
			continue
		}
		if opts.SkipLocked && IsLocked(es, k) {
			continue
		}
		delete(es.Vars, k)
		removed = append(removed, k)
	}
	sort.Strings(removed)
	return removed, nil
}
