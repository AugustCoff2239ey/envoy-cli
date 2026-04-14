package envset

import "fmt"

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchOpSet    PatchOp = "set"
	PatchOpDelete PatchOp = "delete"
)

// PatchEntry describes one operation in a patch.
type PatchEntry struct {
	Op    PatchOp
	Key   string
	Value string
}

// Patch applies a list of patch entries to the given EnvSet.
// Set operations add or update keys; delete operations remove them.
// Returns an error if any entry is invalid or references a locked key.
func Patch(es *EnvSet, entries []PatchEntry) error {
	if es == nil {
		return fmt.Errorf("patch: envset is nil")
	}
	for i, e := range entries {
		if err := ValidateKey(e.Key); err != nil {
			return fmt.Errorf("patch: entry %d: %w", i, err)
		}
		if IsLocked(es, e.Key) {
			return fmt.Errorf("patch: entry %d: key %q is locked", i, e.Key)
		}
		switch e.Op {
		case PatchOpSet:
			if err := es.Set(e.Key, e.Value); err != nil {
				return fmt.Errorf("patch: entry %d: %w", i, err)
			}
		case PatchOpDelete:
			es.Delete(e.Key)
		default:
			return fmt.Errorf("patch: entry %d: unknown op %q", i, e.Op)
		}
	}
	return nil
}

// PatchFromDiff converts a DiffResult into a slice of PatchEntries that, when
// applied to the base set, would make it match the target set.
func PatchFromDiff(d DiffResult) []PatchEntry {
	var entries []PatchEntry
	for k, v := range d.Added {
		entries = append(entries, PatchEntry{Op: PatchOpSet, Key: k, Value: v})
	}
	for k, c := range d.Changed {
		entries = append(entries, PatchEntry{Op: PatchOpSet, Key: k, Value: c.New})
	}
	for k := range d.Removed {
		entries = append(entries, PatchEntry{Op: PatchOpDelete, Key: k})
	}
	return entries
}
