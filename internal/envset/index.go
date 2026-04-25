package envset

import (
	"fmt"
	"sort"
	"strings"
)

// IndexEntry holds positional and metadata info for a key in an EnvSet.
type IndexEntry struct {
	Key      string
	Position int
	Group    string
	Tags     []string
}

// Index builds a lookup map of key -> IndexEntry for fast access and introspection.
func Index(es *EnvSet) (map[string]IndexEntry, error) {
	if es == nil {
		return nil, fmt.Errorf("index: nil EnvSet")
	}

	keys := SortedKeys(es)
	result := make(map[string]IndexEntry, len(keys))

	for i, k := range keys {
		entry := IndexEntry{
			Key:      k,
			Position: i,
		}

		// Attach group if present
		if raw, ok := es.Meta["group:"+k]; ok {
			entry.Group = raw
		}

		// Attach tags if present
		if raw, ok := es.Meta["tags:"+k]; ok && raw != "" {
			parts := strings.Split(raw, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					entry.Tags = append(entry.Tags, p)
				}
			}
		}

		result[k] = entry
	}

	return result, nil
}

// IndexKeys returns a sorted slice of all indexed keys.
func IndexKeys(idx map[string]IndexEntry) []string {
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// IndexByGroup returns all IndexEntry values belonging to a given group.
func IndexByGroup(idx map[string]IndexEntry, group string) []IndexEntry {
	var entries []IndexEntry
	for _, e := range idx {
		if e.Group == group {
			entries = append(entries, e)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Position < entries[j].Position
	})
	return entries
}
