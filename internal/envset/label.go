package envset

import (
	"fmt"
	"regexp"
)

var validLabelKey = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// AddLabel attaches a key=value label to an EnvSet.
func AddLabel(es *EnvSet, key, value string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	if !validLabelKey.MatchString(key) {
		return fmt.Errorf("invalid label key %q: must match [a-zA-Z0-9_\\-.]+", key)
	}
	if es.Meta == nil {
		es.Meta = make(map[string]string)
	}
	es.Meta["label:"+key] = value
	return nil
}

// RemoveLabel removes a label by key from an EnvSet.
func RemoveLabel(es *EnvSet, key string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	mk := "label:" + key
	if _, ok := es.Meta[mk]; !ok {
		return fmt.Errorf("label %q not found", key)
	}
	delete(es.Meta, mk)
	return nil
}

// GetLabel returns the value for a label key.
func GetLabel(es *EnvSet, key string) (string, bool) {
	if es == nil || es.Meta == nil {
		return "", false
	}
	v, ok := es.Meta["label:"+key]
	return v, ok
}

// ListLabels returns all labels as a map.
func ListLabels(es *EnvSet) map[string]string {
	out := make(map[string]string)
	if es == nil || es.Meta == nil {
		return out
	}
	prefix := "label:"
	for k, v := range es.Meta {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out[k[len(prefix):]] = v
		}
	}
	return out
}
