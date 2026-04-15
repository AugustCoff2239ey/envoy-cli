package envset

import (
	"fmt"
	"sort"
	"time"
)

// ArchiveEntry represents a single archived version of an EnvSet.
type ArchiveEntry struct {
	ID        string
	ArchivedAt time.Time
	Reason    string
	Snapshot  map[string]string
}

// Archive holds a collection of archived EnvSet snapshots.
type Archive struct {
	entries []ArchiveEntry
}

// NewArchive creates an empty Archive.
func NewArchive() *Archive {
	return &Archive{}
}

// Add stores a snapshot of the given EnvSet in the archive.
func (a *Archive) Add(es *EnvSet, reason string) (ArchiveEntry, error) {
	if a == nil {
		return ArchiveEntry{}, fmt.Errorf("archive is nil")
	}
	if es == nil {
		return ArchiveEntry{}, fmt.Errorf("envset is nil")
	}
	snap := make(map[string]string, len(es.Vars))
	for k, v := range es.Vars {
		snap[k] = v
	}
	entry := ArchiveEntry{
		ID:        fmt.Sprintf("%s-%s-%d", es.Name, es.Environment, time.Now().UnixNano()),
		ArchivedAt: time.Now().UTC(),
		Reason:    reason,
		Snapshot:  snap,
	}
	a.entries = append(a.entries, entry)
	return entry, nil
}

// List returns all archive entries, sorted newest first.
func (a *Archive) List() []ArchiveEntry {
	if a == nil {
		return nil
	}
	copy := make([]ArchiveEntry, len(a.entries))
	for i, e := range a.entries {
		copy[i] = e
	}
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].ArchivedAt.After(copy[j].ArchivedAt)
	})
	return copy
}

// Get retrieves an archive entry by ID.
func (a *Archive) Get(id string) (ArchiveEntry, bool) {
	for _, e := range a.entries {
		if e.ID == id {
			return e, true
		}
	}
	return ArchiveEntry{}, false
}

// Restore creates a new EnvSet from an archive entry.
func (a *Archive) Restore(id, name, environment string) (*EnvSet, error) {
	entry, ok := a.Get(id)
	if !ok {
		return nil, fmt.Errorf("archive entry %q not found", id)
	}
	es, err := New(name, environment)
	if err != nil {
		return nil, err
	}
	for k, v := range entry.Snapshot {
		if err := es.Set(k, v); err != nil {
			return nil, fmt.Errorf("restore: failed to set %s: %w", k, err)
		}
	}
	return es, nil
}
