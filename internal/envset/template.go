package envset

import (
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches {{VAR_NAME}} placeholders in templates.
var templateVarRe = regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*)\}\}`)

// TemplateResult holds the rendered output and any unresolved placeholders.
type TemplateResult struct {
	Rendered     string
	Unresolved   []string
}

// RenderTemplate replaces {{KEY}} placeholders in tmpl with values from e.
// Keys not found in e are left as-is and recorded in Unresolved.
func RenderTemplate(e *EnvSet, tmpl string) (*TemplateResult, error) {
	if e == nil {
		return nil, fmt.Errorf("envset: RenderTemplate: nil EnvSet")
	}
	if tmpl == "" {
		return &TemplateResult{Rendered: ""}, nil
	}

	unresolved := []string{}
	seen := map[string]bool{}

	rendered := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := templateVarRe.FindStringSubmatch(match)[1]
		val, ok := e.Vars[key]
		if !ok {
			if !seen[key] {
				unresolved = append(unresolved, key)
				seen[key] = true
			}
			return match
		}
		return val
	})

	return &TemplateResult{
		Rendered:   rendered,
		Unresolved: unresolved,
	}, nil
}

// ExtractPlaceholders returns all unique {{KEY}} placeholder names found in tmpl.
func ExtractPlaceholders(tmpl string) []string {
	matches := templateVarRe.FindAllStringSubmatch(tmpl, -1)
	seen := map[string]bool{}
	keys := []string{}
	for _, m := range matches {
		k := m[1]
		if !seen[k] {
			keys = append(keys, k)
			seen[k] = true
		}
	}
	return keys
}

// TemplateComplete returns true if all placeholders in tmpl are satisfied by e.
func TemplateComplete(e *EnvSet, tmpl string) bool {
	if e == nil {
		return false
	}
	for _, key := range ExtractPlaceholders(tmpl) {
		if _, ok := e.Vars[key]; !ok {
			return false
		}
	}
	return true
}

// MissingPlaceholders returns keys referenced in tmpl that are absent from e.
func MissingPlaceholders(e *EnvSet, tmpl string) []string {
	if e == nil {
		return ExtractPlaceholders(tmpl)
	}
	missing := []string{}
	for _, key := range ExtractPlaceholders(tmpl) {
		if _, ok := e.Vars[key]; !ok {
			missing = append(missing, key)
		}
	}
	return missing
}

// SuggestPlaceholders returns keys in e whose names appear as substrings of
// any unresolved placeholder (case-insensitive), useful for typo hints.
func SuggestPlaceholders(e *EnvSet, unresolved []string) map[string][]string {
	suggestions := map[string][]string{}
	for _, u := range unresolved {
		uLower := strings.ToLower(u)
		for k := range e.Vars {
			if strings.Contains(strings.ToLower(k), uLower) || strings.Contains(uLower, strings.ToLower(k)) {
				suggestions[u] = append(suggestions[u], k)
			}
		}
	}
	return suggestions
}
