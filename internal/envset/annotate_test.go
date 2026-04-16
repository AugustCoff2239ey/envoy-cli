package envset

import (
	"testing"
)

func baseAnnotateSet() *EnvSet {
	e, _ := New("annotate-set", "local")
	_ = e.Set("API_KEY", "abc123")
	_ = e.Set("DB_URL", "postgres://localhost/db")
	return e
}

func TestAnnotate_Valid(t *testing.T) {
	e := baseAnnotateSet()
	if err := Annotate(e, "API_KEY", "Primary API key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	note, ok := GetAnnotation(e, "API_KEY")
	if !ok || note != "Primary API key" {
		t.Errorf("expected annotation, got %q", note)
	}
}

func TestAnnotate_NilEnvSet(t *testing.T) {
	if err := Annotate(nil, "API_KEY", "note"); err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestAnnotate_MissingKey(t *testing.T) {
	e := baseAnnotateSet()
	if err := Annotate(e, "MISSING", "note"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestAnnotate_NewlineInNote(t *testing.T) {
	e := baseAnnotateSet()
	if err := Annotate(e, "API_KEY", "bad\nnote"); err == nil {
		t.Error("expected error for newline in note")
	}
}

func TestRemoveAnnotation_Valid(t *testing.T) {
	e := baseAnnotateSet()
	_ = Annotate(e, "DB_URL", "connection string")
	if err := RemoveAnnotation(e, "DB_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := GetAnnotation(e, "DB_URL"); ok {
		t.Error("expected annotation to be removed")
	}
}

func TestRemoveAnnotation_NotFound(t *testing.T) {
	e := baseAnnotateSet()
	if err := RemoveAnnotation(e, "API_KEY"); err == nil {
		t.Error("expected error when annotation not found")
	}
}

func TestListAnnotations_Multiple(t *testing.T) {
	e := baseAnnotateSet()
	_ = Annotate(e, "API_KEY", "note1")
	_ = Annotate(e, "DB_URL", "note2")
	list := ListAnnotations(e)
	if len(list) != 2 {
		t.Errorf("expected 2 annotations, got %d", len(list))
	}
}

func TestListAnnotations_Empty(t *testing.T) {
	e := baseAnnotateSet()
	if list := ListAnnotations(e); list != nil {
		t.Errorf("expected nil, got %v", list)
	}
}
