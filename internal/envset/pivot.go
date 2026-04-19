package envset

import (
	"fmt"
	"sort"
)

// PivotResult holds the result of a pivot operation — grouping envsets by a shared key's value.
type PivotResult struct {
	Key    string
	Groups map[string][]string // value -> list of envset names
}

// Pivot groups multiple EnvSets by the value of a shared key.
// It returns a PivotResult mapping each distinct value to the envset names that hold it.
// If an envset does not contain the key, it is placed under the "" (empty) group.
func Pivot(key string, sets ...*EnvSet) (*PivotResult, error) {
	if err := ValidateKey(key); err != nil {
		return nil, fmt.Errorf("pivot: invalid key %q: %w", key, err)
	}
	if len(sets) == 0 {
		return nil, fmt.Errorf("pivot: no envsets provided")
	}

	result := &PivotResult{
		Key:    key,
		Groups: make(map[string][]string),
	}

	for _, s := range sets {
		if s == nil {
			continue
		}
		val, ok := s.Vars[key]
		if !ok {
			val = ""
		}
		result.Groups[val] = append(result.Groups[val], s.Name)
	}

	// sort names within each group for determinism
	for v := range result.Groups {
		sort.Strings(result.Groups[v])
	}

	return result, nil
}

// PivotKeys returns all distinct values found for the pivot key, sorted.
func (p *PivotResult) PivotKeys() []string {
	keys := make([]string, 0, len(p.Groups))
	for k := range p.Groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
