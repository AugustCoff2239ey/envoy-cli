package envset

import "fmt"

const freezeMetaKey = "__frozen__"

// Freeze marks an EnvSet as frozen, preventing any further modifications.
func Freeze(es *EnvSet) error {
	if es == nil {
		return fmt.Errorf("freeze: envset is nil")
	}
	es.Meta[freezeMetaKey] = "true"
	return nil
}

// Unfreeze removes the frozen state from an EnvSet.
func Unfreeze(es *EnvSet) error {
	if es == nil {
		return fmt.Errorf("unfreeze: envset is nil")
	}
	delete(es.Meta, freezeMetaKey)
	return nil
}

// IsFrozen reports whether the EnvSet is currently frozen.
func IsFrozen(es *EnvSet) bool {
	if es == nil {
		return false
	}
	return es.Meta[freezeMetaKey] == "true"
}

// AssertMutable returns an error if the EnvSet is frozen.
func AssertMutable(es *EnvSet) error {
	if IsFrozen(es) {
		return fmt.Errorf("envset %q (%s) is frozen and cannot be modified", es.Name, es.Environment)
	}
	return nil
}
