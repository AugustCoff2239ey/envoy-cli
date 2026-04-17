package envset

import (
	"os"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	return s
}

func TestStore_SaveAndLoad(t *testing.T) {
	s := tempStore(t)

	es, _ := New("webapp", EnvStaging)
	_ = es.Set("API_URL", "https://staging.example.com")
	_ = es.Set("DEBUG", "true")

	if err := s.Save(es); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := s.Load("webapp", EnvStaging)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != es.Name {
		t.Errorf("expected name '%s', got '%s'", es.Name, loaded.Name)
	}
	if loaded.Environment != es.Environment {
		t.Errorf("expected env '%s', got '%s'", es.Environment, loaded.Environment)
	}
	if v, ok := loaded.Get("API_URL"); !ok || v != "https://staging.example.com" {
		t.Errorf("unexpected value for API_URL: %s", v)
	}
}

func TestStore_LoadNotFound(t *testing.T) {
	s := tempStore(t)
	_, err := s.Load("nonexistent", EnvLocal)
	if err == nil {
		t.Fatal("expected error for missing envset, got nil")
	}
}

func TestStore_Delete(t *testing.T) {
	s := tempStore(t)

	es, _ := New("service", EnvProduction)
	_ = es.Set("SECRET_KEY", "abc123")
	_ = s.Save(es)

	if err := s.Delete("service", EnvProduction); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := s.Load("service", EnvProduction)
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	s := tempStore(t)
	err := s.Delete("ghost", EnvLocal)
	if err == nil {
		t.Fatal("expected error deleting non-existent envset, got nil")
	}
}

func TestStore_SaveNil(t *testing.T) {
	s := tempStore(t)
	if err := s.Save(nil); err == nil {
		t.Fatal("expected error saving nil EnvSet, got nil")
	}
}

func TestNewStore_DefaultDir(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	s, err := NewStore("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.BaseDir == "" {
		t.Error("expected BaseDir to be set")
	}
}

func TestStore_SaveAndOverwrite(t *testing.T) {
	s := tempStore(t)

	es, _ := New("webapp", EnvStaging)
	_ = es.Set("API_URL", "https://staging.example.com")
	_ = s.Save(es)

	// Overwrite with updated value
	_ = es.Set("API_URL", "https://staging-v2.example.com")
	if err := s.Save(es); err != nil {
		t.Fatalf("Save (overwrite) failed: %v", err)
	}

	loaded, err := s.Load("webapp", EnvStaging)
	if err != nil {
		t.Fatalf("Load after overwrite failed: %v", err)
	}
	if v, ok := loaded.Get("API_URL"); !ok || v != "https://staging-v2.example.com" {
		t.Errorf("expected updated API_URL, got %s", v)
	}
}
