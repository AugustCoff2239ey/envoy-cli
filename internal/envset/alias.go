package envset

import (
	"fmt"
	"regexp"
)

var validAliasName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// AddAlias creates a named alias for an existing key.
func AddAlias(es *EnvSet, key, alias string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if !validAliasName.MatchString(alias) {
		return fmt.Errorf("invalid alias name %q: must start with a letter and contain only letters, digits, underscores, or hyphens", alias)
	}
	if _, ok := es.Vars[key]; !ok {
		return fmt.Errorf("key %q does not exist", key)
	}
	if es.Meta == nil {
		es.Meta = make(map[string]map[string]string)
	}
	if es.Meta[key] == nil {
		es.Meta[key] = make(map[string]string)
	}
	es.Meta[key]["alias:"+alias] = key
	return nil
}

// RemoveAlias removes a named alias from a key.
func RemoveAlias(es *EnvSet, key, alias string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if es.Meta == nil || es.Meta[key] == nil {
		return fmt.Errorf("alias %q not found on key %q", alias, key)
	}
	aliasKey := "alias:" + alias
	if _, ok := es.Meta[key][aliasKey]; !ok {
		return fmt.Errorf("alias %q not found on key %q", alias, key)
	}
	delete(es.Meta[key], aliasKey)
	return nil
}

// ListAliases returns all aliases defined for a key.
func ListAliases(es *EnvSet, key string) []string {
	if es == nil || es.Meta == nil || es.Meta[key] == nil {
		return nil
	}
	var aliases []string
	prefix := "alias:"
	for k := range es.Meta[key] {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			aliases = append(aliases, k[len(prefix):])
		}
	}
	return aliases
}

// ResolveAlias returns the key that the given alias points to, if any.
func ResolveAlias(es *EnvSet, alias string) (string, bool) {
	if es == nil || es.Meta == nil {
		return "", false
	}
	aliasKey := "alias:" + alias
	for key, meta := range es.Meta {
		if _, ok := meta[aliasKey]; ok {
			return key, true
		}
	}
	return "", false
}
