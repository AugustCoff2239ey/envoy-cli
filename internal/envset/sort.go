package envset

import (
	"fmt"
	"sort"
	"strings"
)

// SortOptions controls how keys are sorted within an EnvSet.
type SortOptions struct {
	Descending bool
	ByValue    bool
	Keys       []string // if non-empty, sort only these keys (others retain position)
}

// DefaultSortOptions returns the default sort configuration.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Descending: false,
		ByValue:    false,
	}
}

// Sort reorders the keys of an EnvSet according to SortOptions.
// It stores the resulting order in the EnvSet metadata under "_sort_order".
func Sort(es *EnvSet, opts SortOptions) error {
	if es == nil {
		return fmt.Errorf("sort: nil EnvSet")
	}

	keys := make([]string, 0, len(es.Vars))
	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			if _, ok := es.Vars[k]; !ok {
				return fmt.Errorf("sort: key %q not found", k)
			}
			keys = append(keys, k)
		}
	} else {
		for k := range es.Vars {
			keys = append(keys, k)
		}
	}

	if opts.ByValue {
		sort.SliceStable(keys, func(i, j int) bool {
			vi, vj := es.Vars[keys[i]], es.Vars[keys[j]]
			if opts.Descending {
				return vi > vj
			}
			return vi < vj
		})
	} else {
		sort.SliceStable(keys, func(i, j int) bool {
			if opts.Descending {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	}

	if es.Meta == nil {
		es.Meta = map[string]string{}
	}
	es.Meta["_sort_order"] = strings.Join(keys, ",")
	return nil
}

// SortedKeys returns the keys of an EnvSet in the stored sort order,
// falling back to lexicographic order if none is set.
func SortedKeys(es *EnvSet) []string {
	if es == nil {
		return nil
	}
	if es.Meta != nil {
		if order, ok := es.Meta["_sort_order"]; ok && order != "" {
			return strings.Split(order, ",")
		}
	}
	keys := make([]string, 0, len(es.Vars))
	for k := range es.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
