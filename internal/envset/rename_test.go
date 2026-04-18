package envset

import (
	"testing"
)

func baseRenameSet(t *testing.T) (*Store, *EnvSet) {
	t.Helper()
	store := tempStore(t)
	es, err := New("myapp", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("PORT", "5432")
	if err := store.Save(es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return store, es
}

func TestRename_NewName(t *testing.T) {
	store, src := baseRenameSet(t)
	renamed, err := Rename(store, src, RenameOptions{NewName: "webapp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if renamed.Name != "webapp" || renamed.Environment != "staging" {
		t.Errorf("expected webapp/staging, got %s/%s", renamed.Name, renamed.Environment)
	}
	if renamed.Vars["DB_HOST"] != "localhost" {
		t.Errorf("vars not copied correctly")
	}
	// Old key should be gone.
	_, err = store.Load("myapp", "staging")
	if err == nil {
		t.Error("expected old envset to be deleted")
	}
}

func TestRename_NewEnvironment(t *testing.T) {
	store, src := baseRenameSet(t)
	renamed, err := Rename(store, src, RenameOptions{NewEnvironment: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if renamed.Environment != "production" {
		t.Errorf("expected production, got %s", renamed.Environment)
	}
}

func TestRename_NoChange(t *testing.T) {
	store, src := baseRenameSet(t)
	_, err := Rename(store, src, RenameOptions{})
	if err == nil {
		t.Error("expected error when name and env are unchanged")
	}
	_ = store
}

func TestRename_NilSource(t *testing.T) {
	store := tempStore(t)
	_, err := Rename(store, nil, RenameOptions{NewName: "other"})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestRename_InvalidNewName(t *testing.T) {
	store, src := baseRenameSet(t)
	_, err := Rename(store, src, RenameOptions{NewName: "bad name!"})
	if err == nil {
		t.Error("expected error for invalid new name")
	}
}

func TestRename_NewNameAndEnvironment(t *testing.T) {
	store, src := baseRenameSet(t)
	renamed, err := Rename(store, src, RenameOptions{NewName: "webapp", NewEnvironment: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if renamed.Name != "webapp" || renamed.Environment != "production" {
		t.Errorf("expected webapp/production, got %s/%s", renamed.Name, renamed.Environment)
	}
	// Old key should be gone.
	_, err = store.Load("myapp", "staging")
	if err == nil {
		t.Error("expected old envset to be deleted")
	}
}
