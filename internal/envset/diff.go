package envset

import "fmt"

// DiffResult holds the differences between two EnvSets.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [oldVal, newVal]
}

// String returns a human-readable summary of the diff.
func (d *DiffResult) String() string {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return "No differences found."
	}
	out := ""
	for k, v := range d.Added {
		out += fmt.Sprintf("+ %s=%s\n", k, v)
	}
	for k, v := range d.Removed {
		out += fmt.Sprintf("- %s=%s\n", k, v)
	}
	for k, vals := range d.Changed {
		out += fmt.Sprintf("~ %s: %s -> %s\n", k, vals[0], vals[1])
	}
	return out
}

// HasDiff returns true if there are any differences.
func (d *DiffResult) HasDiff() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Diff compares two EnvSets and returns a DiffResult.
// base is the reference (e.g. local), target is what is being compared (e.g. staging).
func Diff(base, target *EnvSet) *DiffResult {
	result := &DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	// Keys in target not in base -> Added
	for k, v := range target.Vars {
		if _, exists := base.Vars[k]; !exists {
			result.Added[k] = v
		}
	}

	// Keys in base not in target -> Removed; keys in both but different -> Changed
	for k, baseVal := range base.Vars {
		if targetVal, exists := target.Vars[k]; !exists {
			result.Removed[k] = baseVal
		} else if baseVal != targetVal {
			result.Changed[k] = [2]string{baseVal, targetVal}
		}
	}

	return result
}
