package envset

import (
	"fmt"
	"regexp"
	"strings"
)

var interpolatePattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(?::([^}]*))?\}`)

// InterpolateResult holds the output of an interpolation pass.
type InterpolateResult struct {
	Resolved   map[string]string
	Unresolved []string
}

// Interpolate expands ${KEY} and ${KEY:default} references within the values
// of e using other keys in the same EnvSet as the source of truth.
// If a referenced key is absent and no default is provided, the placeholder is
// left unchanged and the key name is appended to UnresolvableKeys.
func Interpolate(e *EnvSet) (InterpolateResult, error) {
	if e == nil {
		return InterpolateResult{}, fmt.Errorf("interpolate: nil EnvSet")
	}

	result := InterpolateResult{
		Resolved: make(map[string]string),
	}

	for key, val := range e.Vars {
		expanded, unresolved := expandValue(val, e.Vars)
		result.Resolved[key] = expanded
		result.Unresolved = append(result.Unresolved, unresolved...)
	}

	// deduplicate unresolved
	seen := make(map[string]struct{})
	deduped := result.Unresolved[:0]
	for _, u := range result.Unresolved {
		if _, ok := seen[u]; !ok {
			seen[u] = struct{}{}
			deduped = append(deduped, u)
		}
	}
	result.Unresolved = deduped

	return result, nil
}

func expandValue(val string, vars map[string]string) (string, []string) {
	var unresolved []string
	expanded := interpolatePattern.ReplaceAllStringFunc(val, func(match string) string {
		subs := interpolatePattern.FindStringSubmatch(match)
		if len(subs) < 2 {
			return match
		}
		refKey := subs[1]
		defaultVal := ""
		hasDefault := len(subs) > 2 && subs[2] != ""
		if hasDefault {
			defaultVal = subs[2]
		}
		if v, ok := vars[refKey]; ok {
			return v
		}
		if hasDefault {
			return defaultVal
		}
		unresolved = append(unresolved, refKey)
		return match
	})
	return strings.TrimSpace(expanded), unresolved
}
