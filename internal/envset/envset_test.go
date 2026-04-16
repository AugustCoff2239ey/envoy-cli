package envset

import (
	"testing"
)

func TestNew_ValidInputs(t *testing.T) {
	es, err := New("myapp", EnvLocal)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if es.Name != "myapp" {
		t.Errorf("expected name 'myapp', got '%s'", es.Name)
	}
	if es.Environment != EnvLocal {
		t.Errorf("expected environment 'local', got '%s'", es.Environment)
	}
	if es.Variables == nil {
		t.Error("expected variables map to be initialized")
	}
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("", EnvStaging)
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestNew_InvalidEnvironment(t *testing.T) {
	_, err := New("myapp", Environment("unknown"))
	if err == nil {
		t.Fatal("expected error for invalid environment, got nil")
	}
}

func TestSet_ValidKey(t *testing.T) {
	es, _ := New("myapp", EnvProduction)
	if err := es.Set("DATABASE_URL", "postgres://localhost/db"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	v, ok := es.Get("DATABASE_URL")
	if !ok || v != "postgres://localhost/db" {
		t.Errorf("expected 'postgres://localhost/db', got '%s'", v)
	}
}

func TestSet_InvalidKey(t *testing.T) {
	es, _ := New("myapp", EnvLocal)
	if err := es.Set("invalid-key", "value"); err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if err := es.Set("1STARTS_WITH_DIGIT", "value"); err == nil {
		t.Fatal("expected error for key starting with digit, got nil")
	}
}

func TestDelete(t *testing.T) {
	es, _ := New("myapp", EnvLocal)
	_ = es.Set("API_KEY", "secret")
	es.Delete("API_KEY")
	_, ok := es.Get("API_KEY")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestDelete_NonExistentKey(t *testing.T) {
	// Deleting a key that does not exist should not panic or return an error.
	es, _ := New("myapp", EnvLocal)
	es.Delete("DOES_NOT_EXIST")
	_, ok := es.Get("DOES_NOT_EXIST")
	if ok {
		t.Error("expected key to remain absent after deleting non-existent key")
	}
}

func TestKeys(t *testing.T) {
	es, _ := New("myapp", EnvStaging)
	_ = es.Set("FOO", "1")
	_ = es.Set("BAR", "2")
	keys := es.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestGet_NonExistentKey(t *testing.T) {
	es, _ := New("myapp", EnvLocal)
	v, ok := es.Get("DOES_NOT_EXIST")
	if ok {
		t.Errorf("expected ok=false for missing key, got ok=true with value '%s'", v)
	}
	if v != "" {
		t.Errorf("expected empty string for missing key, got '%s'", v)
	}
}
