package envset

// ErrEmptyQuery is returned when a Search is called with a blank query string.
var ErrEmptyQuery = errorf("search query must not be empty")

// Search scans the key-value pairs of an EnvSet for entries that match the
// given query string. It supports plain substring matching and full regular
// expression patterns, with optional case-insensitive comparison.
//
// By default both keys and values are searched. Use SearchOptions.KeysOnly or
// SearchOptions.ValuesOnly to narrow the scope.
//
// Each SearchResult records which field ("key", "value", or "both") produced
// the match, making it easy for callers to highlight the relevant portion of
// the output.
//
// Example:
//
//	results, err := Search(es, "DATABASE", SearchOptions{KeysOnly: true})
//	for _, r := range results {
//		fmt.Printf("%s = %s\n", r.Key, r.Value)
//	}
