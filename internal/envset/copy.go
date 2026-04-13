package envset

import "fmt"

// CopyOptions controls how variables are copied between EnvSets.
type CopyOptions struct {
	// Overwrite determines whether existing keys in the destination are replaced.
	Overwrite bool
	// Keys restricts the copy to only the specified keys. If empty, all keys are copied.
	Keys []string
}

// Copy copies environment variables from src into dst according to opts.
// It returns the number of keys copied and any error encountered.
func Copy(src, dst *EnvSet, opts CopyOptions) (int, error) {
	if src == nil {
		return 0, fmt.Errorf("copy: source EnvSet must not be nil")
	}
	if dst == nil {
		return 0, fmt.Errorf("copy: destination EnvSet must not be nil")
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src.Vars {
			keys = append(keys, k)
		}
	}

	copied := 0
	for _, k := range keys {
		val, ok := src.Vars[k]
		if !ok {
			return copied, fmt.Errorf("copy: key %q not found in source", k)
		}
		if _, exists := dst.Vars[k]; exists && !opts.Overwrite {
			continue
		}
		if err := dst.Set(k, val); err != nil {
			return copied, fmt.Errorf("copy: failed to set key %q in destination: %w", k, err)
		}
		copied++
	}
	return copied, nil
}
