package envset

import (
	"testing"
)

func baseArchiveSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("myapp", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("API_KEY", "secret123")
	return es
}

func TestArchive_AddAndList(t *testing.T) {
	a := NewArchive()
	es := baseArchiveSet(t)

	entry, err := a.Add(es, "before deploy")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if entry.ID == "" {
		t.Error("expected non-empty ID")
	}
	if entry.Reason != "before deploy" {
		t.Errorf("expected reason 'before deploy', got %q", entry.Reason)
	}

	list := a.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
}

func TestArchive_Get(t *testing.T) {
	a := NewArchive()
	es := baseArchiveSet(t)

	entry, _ := a.Add(es, "checkpoint")

	got, ok := a.Get(entry.ID)
	if !ok {
		t.Fatal("expected to find entry by ID")
	}
	if got.Snapshot["DB_HOST"] != "localhost" {
		t.Errorf("snapshot mismatch: got %q", got.Snapshot["DB_HOST"])
	}

	_, ok = a.Get("nonexistent")
	if ok {
		t.Error("expected not found for unknown ID")
	}
}

func TestArchive_Restore(t *testing.T) {
	a := NewArchive()
	es := baseArchiveSet(t)

	entry, _ := a.Add(es, "pre-update")

	restored, err := a.Restore(entry.ID, "myapp", "production")
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if restored.Name != "myapp" || restored.Environment != "production" {
		t.Errorf("unexpected name/env: %s/%s", restored.Name, restored.Environment)
	}
	val, _ := restored.Get("API_KEY")
	if val != "secret123" {
		t.Errorf("expected API_KEY=secret123, got %q", val)
	}
}

func TestArchive_RestoreNotFound(t *testing.T) {
	a := NewArchive()
	_, err := a.Restore("missing-id", "app", "local")
	if err == nil {
		t.Error("expected error for missing archive entry")
	}
}

func TestArchive_NilReceiver(t *testing.T) {
	var a *Archive
	_, err := a.Add(nil, "reason")
	if err == nil {
		t.Error("expected error for nil archive")
	}
	if list := a.List(); list != nil {
		t.Error("expected nil list from nil archive")
	}
}
