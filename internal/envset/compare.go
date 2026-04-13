package envset

import "fmt"

// CompareResult holds the result of comparing two EnvSets.
type CompareResult struct {
	MatchingKeys   []string
	MismatchedKeys map[string][2]string // key -> [base value, target value]
	OnlyInBase     []string
	OnlyInTarget   []string
	Equal          bool
}

// Compare performs a detailed comparison between two EnvSets and returns
// a CompareResult describing their relationship.
func Compare(base, target *EnvSet) (*CompareResult, error) {
	if base == nil || target == nil {
		return nil, fmt.Errorf("compare: base and target must not be nil")
	}

	result := &CompareResult{
		MismatchedKeys: make(map[string][2]string),
	}

	baseKeys := make(map[string]struct{})
	for k, bv := range base.Vars {
		baseKeys[k] = struct{}{}
		if tv, ok := target.Vars[k]; ok {
			if bv == tv {
				result.MatchingKeys = append(result.MatchingKeys, k)
			} else {
				result.MismatchedKeys[k] = [2]string{bv, tv}
			}
		} else {
			result.OnlyInBase = append(result.OnlyInBase, k)
		}
	}

	for k := range target.Vars {
		if _, exists := baseKeys[k]; !exists {
			result.OnlyInTarget = append(result.OnlyInTarget, k)
		}
	}

	result.Equal = len(result.MismatchedKeys) == 0 &&
		len(result.OnlyInBase) == 0 &&
		len(result.OnlyInTarget) == 0

	return result, nil
}
