package envset

import (
	"testing"
)

func baseVersionSet() *EnvSet {
	es, _ := New("app", "staging")
	es.Vars["DB_HOST"] = "localhost"
	es.Vars["PORT"] = "5432"
	return es
}

func TestSaveVersion_Basic(t *testing.T) {
	vs := NewVersionStore()
	es := baseVersionSet()
	if err := SaveVersion(vs, es, "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels := ListVersions(vs)
	if len(labels) != 1 || labels[0] != "v1" {
		t.Errorf("expected [v1], got %v", labels)
	}
}

func TestSaveVersion_DuplicateLabel(t *testing.T) {
	vs := NewVersionStore()
	es := baseVersionSet()
	_ = SaveVersion(vs, es, "v1")
	if err := SaveVersion(vs, es, "v1"); err == nil {
		t.Error("expected error for duplicate label")
	}
}

func TestSaveVersion_NilInputs(t *testing.T) {
	vs := NewVersionStore()
	es := baseVersionSet()
	if err := SaveVersion(nil, es, "v1"); err == nil {
		t.Error("expected error for nil store")
	}
	if err := SaveVersion(vs, nil, "v1"); err == nil {
		t.Error("expected error for nil envset")
	}
	if err := SaveVersion(vs, es, ""); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestGetVersion_NotFound(t *testing.T) {
	vs := NewVersionStore()
	if _, err := GetVersion(vs, "missing"); err == nil {
		t.Error("expected error for missing version")
	}
}

func TestRestoreVersion_Basic(t *testing.T) {
	vs := NewVersionStore()
	es := baseVersionSet()
	_ = SaveVersion(vs, es, "v1")

	es.Vars["DB_HOST"] = "changed"
	es.Vars["NEW_KEY"] = "new"

	if err := RestoreVersion(vs, es, "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", es.Vars["DB_HOST"])
	}
}

func TestRestoreVersion_Independence(t *testing.T) {
	vs := NewVersionStore()
	es := baseVersionSet()
	_ = SaveVersion(vs, es, "v1")

	v, _ := GetVersion(vs, "v1")
	v.Vars["DB_HOST"] = "tampered"

	es2, _ := New("app", "staging")
	_ = RestoreVersion(vs, es2, "v1")
	if es2.Vars["DB_HOST"] == "tampered" {
		t.Error("version store was mutated externally")
	}
}
