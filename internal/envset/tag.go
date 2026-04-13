package envset

import (
	"fmt"
	"regexp"
	"sort"
)

var validTagRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// Tag represents a label that can be attached to an EnvSet for grouping or filtering.
type Tag struct {
	Name string
}

// NewTag creates a new Tag, validating the name.
func NewTag(name string) (Tag, error) {
	if name == "" {
		return Tag{}, fmt.Errorf("tag name must not be empty")
	}
	if !validTagRe.MatchString(name) {
		return Tag{}, fmt.Errorf("tag name %q contains invalid characters: only alphanumeric, underscore, and hyphen allowed", name)
	}
	return Tag{Name: name}, nil
}

// AddTag attaches a tag to an EnvSet, returning an error if the tag already exists.
func AddTag(es *EnvSet, tag Tag) error {
	if es == nil {
		return fmt.Errorf("envset must not be nil")
	}
	for _, t := range es.Tags {
		if t == tag.Name {
			return fmt.Errorf("tag %q already exists on envset %q", tag.Name, es.Name)
		}
	}
	es.Tags = append(es.Tags, tag.Name)
	sort.Strings(es.Tags)
	return nil
}

// RemoveTag detaches a tag from an EnvSet, returning an error if not found.
func RemoveTag(es *EnvSet, tagName string) error {
	if es == nil {
		return fmt.Errorf("envset must not be nil")
	}
	for i, t := range es.Tags {
		if t == tagName {
			es.Tags = append(es.Tags[:i], es.Tags[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("tag %q not found on envset %q", tagName, es.Name)
}

// HasTag reports whether an EnvSet has the given tag.
func HasTag(es *EnvSet, tagName string) bool {
	for _, t := range es.Tags {
		if t == tagName {
			return true
		}
	}
	return false
}

// FilterByTag returns a slice of EnvSets from the store that carry the given tag.
func FilterByTag(store *Store, tagName string) ([]*EnvSet, error) {
	all, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("listing envsets: %w", err)
	}
	var result []*EnvSet
	for _, es := range all {
		if HasTag(es, tagName) {
			result = append(result, es)
		}
	}
	return result, nil
}
