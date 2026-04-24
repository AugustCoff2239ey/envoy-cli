package envset

import (
	"fmt"
	"sort"
)

// CrossRefResult holds the result of a cross-reference analysis between two EnvSets.
type CrossRefResult struct {
	// SharedKeys are keys present in both sets.
	SharedKeys []string
	// OnlyInBase are keys only in the base set.
	OnlyInBase []string
	// OnlyInTarget are keys only in the target set.
	OnlyInTarget []string
	// ValueMatches are shared keys where values are identical.
	ValueMatches []string
	// ValueMismatches are shared keys where values differ.
	ValueMismatches []string
}

// CrossRef performs a detailed cross-reference between base and target EnvSets,
// identifying shared keys, exclusive keys, and value agreement.
func CrossRef(base, target *EnvSet) (*CrossRefResult, error) {
	if base == nil || target == nil {
		return nil, fmt.Errorf("crossref: base and target must not be nil")
	}

	result := &CrossRefResult{}

	baseKeys := make(map[string]string)
	for k, v := range base.Vars {
		baseKeys[k] = v
	}

	targetKeys := make(map[string]string)
	for k, v := range target.Vars {
		targetKeys[k] = v
	}

	for k, bv := range baseKeys {
		if tv, ok := targetKeys[k]; ok {
			result.SharedKeys = append(result.SharedKeys, k)
			if bv == tv {
				result.ValueMatches = append(result.ValueMatches, k)
			} else {
				result.ValueMismatches = append(result.ValueMismatches, k)
			}
		} else {
			result.OnlyInBase = append(result.OnlyInBase, k)
		}
	}

	for k := range targetKeys {
		if _, ok := baseKeys[k]; !ok {
			result.OnlyInTarget = append(result.OnlyInTarget, k)
		}
	}

	sort.Strings(result.SharedKeys)
	sort.Strings(result.OnlyInBase)
	sort.Strings(result.OnlyInTarget)
	sort.Strings(result.ValueMatches)
	sort.Strings(result.ValueMismatches)

	return result, nil
}
