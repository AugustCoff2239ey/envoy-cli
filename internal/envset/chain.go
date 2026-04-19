package envset

import "fmt"

// ChainResult holds the result of a chained apply across multiple EnvSets.
type ChainResult struct {
	Applied []string
	Skipped []string
}

// Chain applies a transformation function sequentially across a slice of EnvSets.
// If stopOnError is true, the chain halts on the first error.
func Chain(sets []*EnvSet, fn func(*EnvSet) error, stopOnError bool) (*ChainResult, error) {
	if fn == nil {
		return nil, fmt.Errorf("chain: transform function must not be nil")
	}

	result := &ChainResult{}

	for _, es := range sets {
		if es == nil {
			result.Skipped = append(result.Skipped, "<nil>")
			continue
		}
		if err := fn(es); err != nil {
			result.Skipped = append(result.Skipped, es.Name)
			if stopOnError {
				return result, fmt.Errorf("chain: error on %q: %w", es.Name, err)
			}
		} else {
			result.Applied = append(result.Applied, es.Name)
		}
	}

	return result, nil
}
