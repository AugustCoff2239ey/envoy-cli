package envset

import (
	"fmt"
	"time"
)

// Version represents a named, timestamped snapshot of an EnvSet.
type Version struct {
	Label     string
	CreatedAt time.Time
	Vars      map[string]string
}

// VersionStore holds named versions of an EnvSet.
type VersionStore struct {
	versions []Version
}

// NewVersionStore creates an empty VersionStore.
func NewVersionStore() *VersionStore {
	return &VersionStore{}
}

// SaveVersion saves the current state of an EnvSet under a label.
func SaveVersion(vs *VersionStore, es *EnvSet, label string) error {
	if vs == nil {
		return fmt.Errorf("version store is nil")
	}
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if label == "" {
		return fmt.Errorf("version label must not be empty")
	}
	for _, v := range vs.versions {
		if v.Label == label {
			return fmt.Errorf("version %q already exists", label)
		}
	}
	copy := make(map[string]string, len(es.Vars))
	for k, v := range es.Vars {
		copy[k] = v
	}
	vs.versions = append(vs.versions, Version{
		Label:     label,
		CreatedAt: time.Now(),
		Vars:      copy,
	})
	return nil
}

// GetVersion retrieves a named version.
func GetVersion(vs *VersionStore, label string) (*Version, error) {
	if vs == nil {
		return nil, fmt.Errorf("version store is nil")
	}
	for i := range vs.versions {
		if vs.versions[i].Label == label {
			return &vs.versions[i], nil
		}
	}
	return nil, fmt.Errorf("version %q not found", label)
}

// ListVersions returns all version labels in order.
func ListVersions(vs *VersionStore) []string {
	if vs == nil {
		return nil
	}
	labels := make([]string, len(vs.versions))
	for i, v := range vs.versions {
		labels[i] = v.Label
	}
	return labels
}

// RestoreVersion applies a named version's vars back onto an EnvSet.
func RestoreVersion(vs *VersionStore, es *EnvSet, label string) error {
	v, err := GetVersion(vs, label)
	if err != nil {
		return err
	}
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	for k, val := range v.Vars {
		es.Vars[k] = val
	}
	return nil
}
