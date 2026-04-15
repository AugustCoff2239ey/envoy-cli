package envset

import (
	"fmt"
	"sort"
	"strings"
)

// Group represents a named collection of keys within an EnvSet.
type Group struct {
	Name string
	Keys []string
}

// CreateGroup creates a named group of keys in the EnvSet metadata.
// Returns an error if any key does not exist or the group name is invalid.
func CreateGroup(es *EnvSet, name string, keys []string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if name == "" {
		return fmt.Errorf("group name must not be empty")
	}
	if !isValidGroupName(name) {
		return fmt.Errorf("invalid group name %q: use only letters, digits, hyphens, underscores", name)
	}
	for _, k := range keys {
		if _, ok := es.Vars[k]; !ok {
			return fmt.Errorf("key %q does not exist in envset", k)
		}
	}
	if es.Meta == nil {
		es.Meta = make(map[string]string)
	}
	es.Meta["group:"+name] = strings.Join(keys, ",")
	return nil
}

// GetGroup returns the Group for the given name, or an error if not found.
func GetGroup(es *EnvSet, name string) (*Group, error) {
	if es == nil {
		return nil, fmt.Errorf("envset is nil")
	}
	val, ok := es.Meta["group:"+name]
	if !ok {
		return nil, fmt.Errorf("group %q not found", name)
	}
	keys := []string{}
	if val != "" {
		keys = strings.Split(val, ",")
	}
	return &Group{Name: name, Keys: keys}, nil
}

// ListGroups returns all groups defined in the EnvSet, sorted by name.
func ListGroups(es *EnvSet) []Group {
	if es == nil {
		return nil
	}
	var groups []Group
	for k, v := range es.Meta {
		if strings.HasPrefix(k, "group:") {
			name := strings.TrimPrefix(k, "group:")
			keys := []string{}
			if v != "" {
				keys = strings.Split(v, ",")
			}
			groups = append(groups, Group{Name: name, Keys: keys})
		}
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// DeleteGroup removes a group definition from the EnvSet.
func DeleteGroup(es *EnvSet, name string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	key := "group:" + name
	if _, ok := es.Meta[key]; !ok {
		return fmt.Errorf("group %q not found", name)
	}
	delete(es.Meta, key)
	return nil
}

func isValidGroupName(name string) bool {
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}
