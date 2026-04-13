package envset

import (
	"testing"
)

func baseSnapshotSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("myapp", "staging")
	if err != nil {
		t.Fatalf("failed to create base EnvSet: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	return es
}

func TestTakeSnapshot_Basic(t *testing.T) {
	es := baseSnapshotSet(t)
	snap, err := TakeSnapshot(es, "before migration")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Name != es.Name {
		t.Errorf("expected name %q, got %q", es.Name, snap.Name)
	}
	if snap.Env != es.Environment {
		t.Errorf("expected env %q, got %q", es.Environment, snap.Env)
	}
	if snap.Message != "before migration" {
		t.Errorf("expected message %q, got %q", "before migration", snap.Message)
	}
	if len(snap.Vars) != len(es.Vars) {
		t.Errorf("expected %d vars, got %d", len(es.Vars), len(snap.Vars))
	}
}

func TestTakeSnapshot_Independence(t *testing.T) {
	es := baseSnapshotSet(t)
	snap, _ := TakeSnapshot(es, "")
	_ = es.Set("NEW_KEY", "newval")
	if _, ok := snap.Vars["NEW_KEY"]; ok {
		t.Error("snapshot should not reflect mutations to original EnvSet")
	}
}

func TestTakeSnapshot_NilSource(t *testing.T) {
	_, err := TakeSnapshot(nil, "")
	if err == nil {
		t.Error("expected error for nil source, got nil")
	}
}

func TestRestoreSnapshot_Basic(t *testing.T) {
	es := baseSnapshotSet(t)
	snap, _ := TakeSnapshot(es, "restore test")
	restored, err := RestoreSnapshot(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if restored.Name != snap.Name {
		t.Errorf("expected name %q, got %q", snap.Name, restored.Name)
	}
	for k, v := range snap.Vars {
		if got := restored.Vars[k]; got != v {
			t.Errorf("key %q: expected %q, got %q", k, v, got)
		}
	}
}

func TestRestoreSnapshot_Nil(t *testing.T) {
	_, err := RestoreSnapshot(nil)
	if err == nil {
		t.Error("expected error for nil snapshot, got nil")
	}
}
