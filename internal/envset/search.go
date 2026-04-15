package envset

import (
	"regexp"
	"strings"
)

// SearchOptions controls how Search behaves.
type SearchOptions struct {
	CaseSensitive bool
	Regex         bool
	KeysOnly      bool
	ValuesOnly    bool
}

// SearchResult holds a single match from a Search operation.
type SearchResult struct {
	Key   string
	Value string
	Field string // "key", "value", or "both"
}

// Search scans an EnvSet for entries whose key or value match the given query.
// It respects SearchOptions for case sensitivity, regex mode, and field scope.
func Search(es *EnvSet, query string, opts SearchOptions) ([]SearchResult, error) {
	if es == nil {
		return nil, ErrNilEnvSet
	}
	if query == "" {
		return nil, ErrEmptyQuery
	}

	var matchFn func(s string) bool

	if opts.Regex {
		flags := ""
		if !opts.CaseSensitive {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + query)
		if err != nil {
			return nil, err
		}
		matchFn = re.MatchString
	} else {
		needle := query
		if !opts.CaseSensitive {
			needle = strings.ToLower(query)
		}
		matchFn = func(s string) bool {
			if !opts.CaseSensitive {
				s = strings.ToLower(s)
			}
			return strings.Contains(s, needle)
		}
	}

	var results []SearchResult
	for _, k := range sortedKeys(es.Vars) {
		v := es.Vars[k]
		keyMatch := !opts.ValuesOnly && matchFn(k)
		valMatch := !opts.KeysOnly && matchFn(v)

		if keyMatch || valMatch {
			field := "both"
			if keyMatch && !valMatch {
				field = "key"
			} else if valMatch && !keyMatch {
				field = "value"
			}
			results = append(results, SearchResult{Key: k, Value: v, Field: field})
		}
	}
	return results, nil
}
