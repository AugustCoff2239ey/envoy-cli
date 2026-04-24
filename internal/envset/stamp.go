package envset

import (
	"fmt"
	"time"
)

// StampOptions controls how timestamps are applied to an EnvSet.
type StampOptions struct {
	// Keys to stamp; if empty, all keys are stamped.
	Keys []string
	// Format is a Go time layout string. Defaults to time.RFC3339.
	Format string
	// Prefix is prepended to the generated timestamp value.
	Prefix string
	// Suffix is appended to the generated timestamp value.
	Suffix string
}

// DefaultStampOptions returns sensible defaults for StampOptions.
func DefaultStampOptions() StampOptions {
	return StampOptions{
		Format: time.RFC3339,
	}
}

// Stamp writes the current UTC timestamp as the value for the specified keys
// (or all keys when none are specified). Locked or protected keys are skipped.
// Returns a map of key → stamped value for every key that was updated.
func Stamp(es *EnvSet, opts StampOptions) (map[string]string, error) {
	if es == nil {
		return nil, fmt.Errorf("stamp: nil EnvSet")
	}
	if err := AssertMutable(es); err != nil {
		return nil, fmt.Errorf("stamp: %w", err)
	}
	if err := AssertWritable(es); err != nil {
		return nil, fmt.Errorf("stamp: %w", err)
	}

	format := opts.Format
	if format == "" {
		format = time.RFC3339
	}

	targets := opts.Keys
	if len(targets) == 0 {
		for k := range es.Vars {
			targets = append(targets, k)
		}
	}

	stamped := make(map[string]string)
	now := time.Now().UTC()

	for _, k := range targets {
		if _, exists := es.Vars[k]; !exists {
			return nil, fmt.Errorf("stamp: key %q not found", k)
		}
		if IsLocked(es, k) || IsProtected(es, k) {
			continue
		}
		val := opts.Prefix + now.Format(format) + opts.Suffix
		es.Vars[k] = val
		stamped[k] = val
	}

	return stamped, nil
}
