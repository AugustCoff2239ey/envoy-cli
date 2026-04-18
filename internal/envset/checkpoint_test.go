package envset

import (
	"testing"
)

func baseCheckpointSet() *EnvSet {
	e, _ := New("app", "local")
	_ = e.Set("DB_HOST", "localhost")
	_ = e.Set("DB_PORT", "5432")
	return e
}

func TestCheckpoint_SaveAndLoad(t *testing.T) {
	cs := NewCheckpointStore()
	e := baseCheckpointSet()
	if err := cs.Save("v1", e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	cp, err := cs.Load("v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cp.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", cp.Vars["DB_HOST"])
	}
}

func TestCheckpoint_Independence(t *testing.T) {
	cs := NewCheckpointStore()
	e := baseCheckpointSet()
	_ = cs.Save("v1", e)
	e.Vars["DB_HOST"] = "remotehost"
	cp, _ := cs.Load("v1")
	if cp.Vars["DB_HOST"] != "localhost" {
		t.Error("checkpoint should be independent of source EnvSet")
	}
}

func TestCheckpoint_Restore(t *testing.T) {
	cs := NewCheckpointStore()
	e := baseCheckpointSet()
	_ = cs.Save("v1", e)
	e.Vars["DB_HOST"] = "changed"
	if err := cs.Restore("v1", e); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if e.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost after restore, got %s", e.Vars["DB_HOST"])
	}
}

func TestCheckpoint_LoadNotFound(t *testing.T) {
	cs := NewCheckpointStore()
	_, err := cs.Load("missing")
	if err == nil {
		t.Error("expected error for missing checkpoint")
	}
}

func TestCheckpoint_Delete(t *testing.T) {
	cs := NewCheckpointStore()
	e := baseCheckpointSet()
	_ = cs.Save("v1", e)
	if err := cs.Delete("v1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := cs.Delete("v1"); err == nil {
		t.Error("expected error deleting already-deleted checkpoint")
	}
}

func TestCheckpoint_NilEnvSet(t *testing.T) {
	cs := NewCheckpointStore()
	if err := cs.Save("v1", nil); err == nil {
		t.Error("expected error saving nil EnvSet")
	}
	if err := cs.Restore("v1", nil); err == nil {
		t.Error("expected error restoring into nil EnvSet")
	}
}

func TestCheckpoint_EmptyName(t *testing.T) {
	cs := NewCheckpointStore()
	e := baseCheckpointSet()
	if err := cs.Save("", e); err == nil {
		t.Error("expected error for empty checkpoint name")
	}
}
