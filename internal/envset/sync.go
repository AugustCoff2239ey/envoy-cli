package envset

import "fmt"

// SyncMode controls how keys are merged during a sync.
type SyncMode int

const (
	// SyncModeAddOnly only copies keys from src that are missing in dst.
	SyncModeAddOnly SyncMode = iota
	// SyncModeOverwrite copies all keys from src into dst, overwriting existing.
	SyncModeOverwrite
)

// SyncResult summarises what changed after a Sync call.
type SyncResult struct {
	Added   []string
	Updated []string
}

// String returns a human-readable summary of the sync result.
func (r *SyncResult) String() string {
	if len(r.Added) == 0 && len(r.Updated) == 0 {
		return "Nothing to sync."
	}
	out := fmt.Sprintf("Added: %d, Updated: %d\n", len(r.Added), len(r.Updated))
	for _, k := range r.Added {
		out += fmt.Sprintf("  + %s\n", k)
	}
	for _, k := range r.Updated {
		out += fmt.Sprintf("  ~ %s\n", k)
	}
	return out
}

// Sync copies variables from src into dst according to the given SyncMode.
// It returns a SyncResult describing what changed.
func Sync(src, dst *EnvSet, mode SyncMode) (*SyncResult, error) {
	if src == nil || dst == nil {
		return nil, fmt.Errorf("sync: src and dst must not be nil")
	}

	result := &SyncResult{}

	for k, v := range src.Vars {
		if existing, exists := dst.Vars[k]; !exists {
			dst.Vars[k] = v
			result.Added = append(result.Added, k)
		} else if mode == SyncModeOverwrite && existing != v {
			dst.Vars[k] = v
			result.Updated = append(result.Updated, k)
		}
	}

	return result, nil
}
