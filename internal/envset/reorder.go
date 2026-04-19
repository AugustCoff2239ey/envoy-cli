package envset

import "fmt"

// ReorderResult holds the outcome of a Reorder operation.
type ReorderResult struct {
	Keys []string
}

// Reorder rearranges the display/export order of keys in an EnvSet by storing
// an explicit ordering annotation. Keys not listed retain their relative order
// after the listed ones.
func Reorder(es *EnvSet, order []string) (*ReorderResult, error) {
	if es == nil {
		return nil, fmt.Errorf("reorder: nil EnvSet")
	}
	if len(order) == 0 {
		return nil, fmt.Errorf("reorder: order list is empty")
	}

	// Validate all requested keys exist.
	for _, k := range order {
		if _, ok := es.Vars[k]; !ok {
			return nil, fmt.Errorf("reorder: key %q not found", k)
		}
	}

	// Build ordered list: explicit first, then remaining in original order.
	seen := make(map[string]bool, len(order))
	result := make([]string, 0, len(es.Vars))
	for _, k := range order {
		if !seen[k] {
			result = append(result, k)
			seen[k] = true
		}
	}
	for _, k := range sortedKeys(es.Vars) {
		if !seen[k] {
			result = append(result, k)
		}
	}

	// Persist order as a metadata annotation.
	if es.Meta == nil {
		es.Meta = make(map[string]string)
	}
	es.Meta["key_order"] = joinKeys(result)

	return &ReorderResult{Keys: result}, nil
}

// GetOrder returns the stored key order for an EnvSet, or nil if none is set.
func GetOrder(es *EnvSet) []string {
	if es == nil || es.Meta == nil {
		return nil
	}
	v, ok := es.Meta["key_order"]
	if !ok || v == "" {
		return nil
	}
	return splitKeys(v)
}

func joinKeys(keys []string) string {
	out := ""
	for i, k := range keys {
		if i > 0 {
			out += ","
		}
		out += k
	}
	return out
}

func splitKeys(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	cur := ""
	for _, c := range s {
		if c == ',' {
			out = append(out, cur)
			cur = ""
		} else {
			cur += string(c)
		}
	}
	out = append(out, cur)
	return out
}
