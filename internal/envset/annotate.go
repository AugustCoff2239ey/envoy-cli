package envset

import (
	"fmt"
	"strings"
)

// Annotation holds a note attached to a specific key.
type Annotation struct {
	Key  string
	Note string
}

// Annotate sets a note on the given key in the EnvSet.
func Annotate(e *EnvSet, key, note string) error {
	if e == nil {
		return fmt.Errorf("annotate: nil EnvSet")
	}
	if err := ValidateKey(key); err != nil {
		return fmt.Errorf("annotate: %w", err)
	}
	if _, ok := e.Vars[key]; !ok {
		return fmt.Errorf("annotate: key %q not found", key)
	}
	if strings.Contains(note, "\n") {
		return fmt.Errorf("annotate: note must not contain newlines")
	}
	if e.Annotations == nil {
		e.Annotations = make(map[string]string)
	}
	e.Annotations[key] = note
	return nil
}

// RemoveAnnotation removes the note from the given key.
func RemoveAnnotation(e *EnvSet, key string) error {
	if e == nil {
		return fmt.Errorf("remove annotation: nil EnvSet")
	}
	if e.Annotations == nil {
		return fmt.Errorf("remove annotation: no annotations present")
	}
	if _, ok := e.Annotations[key]; !ok {
		return fmt.Errorf("remove annotation: key %q has no annotation", key)
	}
	delete(e.Annotations, key)
	return nil
}

// GetAnnotation returns the note for the given key.
func GetAnnotation(e *EnvSet, key string) (string, bool) {
	if e == nil || e.Annotations == nil {
		return "", false
	}
	note, ok := e.Annotations[key]
	return note, ok
}

// ListAnnotations returns all annotations as a slice.
func ListAnnotations(e *EnvSet) []Annotation {
	if e == nil || len(e.Annotations) == 0 {
		return nil
	}
	out := make([]Annotation, 0, len(e.Annotations))
	for k, n := range e.Annotations {
		out = append(out, Annotation{Key: k, Note: n})
	}
	return out
}
