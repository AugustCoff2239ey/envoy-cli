package envset

import "fmt"

const metaReadonlyPrefix = "__readonly__"

// MarkReadonly marks an EnvSet as read-only, preventing any modifications.
func MarkReadonly(es *EnvSet) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	es.Meta[metaReadonlyPrefix+"enabled"] = "true"
	return nil
}

// UnmarkReadonly removes the read-only flag from an EnvSet.
func UnmarkReadonly(es *EnvSet) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	delete(es.Meta, metaReadonlyPrefix+"enabled")
	return nil
}

// IsReadonly returns true if the EnvSet is marked as read-only.
func IsReadonly(es *EnvSet) bool {
	if es == nil {
		return false
	}
	return es.Meta[metaReadonlyPrefix+"enabled"] == "true"
}

// AssertWritable returns an error if the EnvSet is read-only.
func AssertWritable(es *EnvSet) error {
	if IsReadonly(es) {
		return fmt.Errorf("envset %q (%s) is read-only", es.Name, es.Environment)
	}
	return nil
}
