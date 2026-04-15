// Package envset provides the Archive type for storing and restoring
// historical snapshots of EnvSets.
//
// An Archive maintains an ordered list of ArchiveEntry values, each
// capturing the full variable map of an EnvSet at a point in time along
// with a human-readable reason string (e.g. "before deploy").
//
// Usage:
//
//	a := envset.NewArchive()
//	entry, err := a.Add(mySet, "pre-release checkpoint")
//
//	// Later, restore from a specific entry:
//	restored, err := a.Restore(entry.ID, "myapp", "production")
//
// Archive entries are immutable snapshots; changes to the original EnvSet
// after archiving do not affect stored entries.
package envset
