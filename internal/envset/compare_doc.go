package envset

// Compare performs a key-by-key comparison between two EnvSets and returns
// a CompareResult describing which keys match, differ, or are exclusive to
// either set.
//
// Usage:
//
//	result, err := Compare(base, target)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, k := range result.Matching {
//	    fmt.Println("match:", k)
//	}
//	for _, k := range result.Mismatched {
//	    fmt.Println("mismatch:", k)
//	}
//
// Fields:
//   - Matching:    keys present in both sets with equal values
//   - Mismatched:  keys present in both sets with differing values
//   - OnlyInBase:  keys found only in the base set
//   - OnlyInTarget: keys found only in the target set
var _ = Compare // ensure Compare is referenced
