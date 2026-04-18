package envset

import (
	"fmt"
	"strings"
)

// DependencyMap maps a key to the list of keys it depends on.
type DependencyMap map[string][]string

// AddDependency records that key depends on dep within the EnvSet.
func AddDependency(es *EnvSet, key, dep string) error {
	if es == nil {
		return fmt.Errorf("envset: nil EnvSet")
	}
	if err := ValidateKey(key); err != nil {
		return fmt.Errorf("envset: invalid key %q: %w", key, err)
	}
	if err := ValidateKey(dep); err != nil {
		return fmt.Errorf("envset: invalid dep %q: %w", dep, err)
	}
	if _, ok := es.Vars[key]; !ok {
		return fmt.Errorf("envset: key %q not found", key)
	}
	if _, ok := es.Vars[dep]; !ok {
		return fmt.Errorf("envset: dep key %q not found", dep)
	}
	if es.Meta == nil {
		es.Meta = map[string]string{}
	}
	metaKey := "dep:" + key
	existing := es.Meta[metaKey]
	deps := splitDeps(existing)
	for _, d := range deps {
		if d == dep {
			return nil
		}
	}
	deps = append(deps, dep)
	es.Meta[metaKey] = strings.Join(deps, ",")
	return nil
}

// RemoveDependency removes dep from key's dependency list.
func RemoveDependency(es *EnvSet, key, dep string) error {
	if es == nil {
		return fmt.Errorf("envset: nil EnvSet")
	}
	metaKey := "dep:" + key
	deps := splitDeps(es.Meta[metaKey])
	filtered := deps[:0]
	for _, d := range deps {
		if d != dep {
			filtered = append(filtered, d)
		}
	}
	if es.Meta == nil {
		es.Meta = map[string]string{}
	}
	es.Meta[metaKey] = strings.Join(filtered, ",")
	return nil
}

// GetDependencies returns the list of keys that key depends on.
func GetDependencies(es *EnvSet, key string) []string {
	if es == nil || es.Meta == nil {
		return nil
	}
	return splitDeps(es.Meta["dep:"+key])
}

// CheckDependencies verifies all declared dependencies for every key are present.
func CheckDependencies(es *EnvSet) []string {
	var missing []string
	if es == nil {
		return missing
	}
	for metaKey, val := range es.Meta {
		if !strings.HasPrefix(metaKey, "dep:") {
			continue
		}
		for _, dep := range splitDeps(val) {
			if _, ok := es.Vars[dep]; !ok {
				missing = append(missing, dep)
			}
		}
	}
	return missing
}

func splitDeps(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
