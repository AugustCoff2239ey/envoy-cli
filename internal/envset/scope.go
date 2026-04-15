package envset

import (
	"fmt"
	"regexp"
)

var validScopeNameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Scope represents a named namespace that groups keys within an EnvSet.
type Scope struct {
	Name string
	Keys []string
}

// CreateScope adds a new scope containing the given keys to the EnvSet.
// All keys must exist in the EnvSet. Returns an error if the scope already
// exists, the name is invalid, or any key is missing.
func CreateScope(es *EnvSet, name string, keys []string) error {
	if es == nil {
		return fmt.Errorf("envset: nil EnvSet")
	}
	if name == "" {
		return fmt.Errorf("scope: name must not be empty")
	}
	if !validScopeNameRe.MatchString(name) {
		return fmt.Errorf("scope: invalid name %q (must start with a letter, alphanumeric/underscore/hyphen only)", name)
	}
	if _, exists := es.Scopes[name]; exists {
		return fmt.Errorf("scope: %q already exists", name)
	}
	for _, k := range keys {
		if _, ok := es.Vars[k]; !ok {
			return fmt.Errorf("scope: key %q not found in EnvSet", k)
		}
	}
	if es.Scopes == nil {
		es.Scopes = make(map[string][]string)
	}
	copy := make([]string, len(keys))
	for i, k := range keys {
		copy[i] = k
	}
	es.Scopes[name] = copy
	return nil
}

// GetScope returns the keys belonging to the named scope.
func GetScope(es *EnvSet, name string) ([]string, error) {
	if es == nil {
		return nil, fmt.Errorf("envset: nil EnvSet")
	}
	keys, ok := es.Scopes[name]
	if !ok {
		return nil, fmt.Errorf("scope: %q not found", name)
	}
	out := make([]string, len(keys))
	copy(out, keys)
	return out, nil
}

// ListScopes returns the names of all scopes defined in the EnvSet.
func ListScopes(es *EnvSet) []string {
	if es == nil || len(es.Scopes) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(es.Scopes))
	for n := range es.Scopes {
		names = append(names, n)
	}
	return names
}

// DeleteScope removes a scope by name. Returns an error if not found.
func DeleteScope(es *EnvSet, name string) error {
	if es == nil {
		return fmt.Errorf("envset: nil EnvSet")
	}
	if _, ok := es.Scopes[name]; !ok {
		return fmt.Errorf("scope: %q not found", name)
	}
	delete(es.Scopes, name)
	return nil
}

// ScopeVars returns a map of key→value for all keys in the named scope.
func ScopeVars(es *EnvSet, name string) (map[string]string, error) {
	keys, err := GetScope(es, name)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = es.Vars[k]
	}
	return out, nil
}
