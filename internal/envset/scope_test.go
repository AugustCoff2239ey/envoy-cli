package envset

import (
	"testing"
)

func baseScopeSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("scope-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, kv := range []struct{ k, v string }{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"API_KEY", "secret"},
		{"LOG_LEVEL", "debug"},
	} {
		if err := es.Set(kv.k, kv.v); err != nil {
			t.Fatalf("Set %s: %v", kv.k, err)
		}
	}
	return es
}

func TestCreateScope_Valid(t *testing.T) {
	es := baseScopeSet(t)
	if err := CreateScope(es, "database", []string{"DB_HOST", "DB_PORT"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, err := GetScope(es, "database")
	if err != nil {
		t.Fatalf("GetScope: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestCreateScope_EmptyName(t *testing.T) {
	es := baseScopeSet(t)
	if err := CreateScope(es, "", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestCreateScope_InvalidName(t *testing.T) {
	es := baseScopeSet(t)
	if err := CreateScope(es, "123bad", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for invalid name")
	}
}

func TestCreateScope_MissingKey(t *testing.T) {
	es := baseScopeSet(t)
	if err := CreateScope(es, "bad", []string{"MISSING_KEY"}); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestCreateScope_Duplicate(t *testing.T) {
	es := baseScopeSet(t)
	_ = CreateScope(es, "database", []string{"DB_HOST"})
	if err := CreateScope(es, "database", []string{"DB_PORT"}); err == nil {
		t.Error("expected error for duplicate scope")
	}
}

func TestDeleteScope_Valid(t *testing.T) {
	es := baseScopeSet(t)
	_ = CreateScope(es, "database", []string{"DB_HOST"})
	if err := DeleteScope(es, "database"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := GetScope(es, "database"); err == nil {
		t.Error("expected error after deletion")
	}
}

func TestDeleteScope_NotFound(t *testing.T) {
	es := baseScopeSet(t)
	if err := DeleteScope(es, "ghost"); err == nil {
		t.Error("expected error for non-existent scope")
	}
}

func TestScopeVars_Valid(t *testing.T) {
	es := baseScopeSet(t)
	_ = CreateScope(es, "db", []string{"DB_HOST", "DB_PORT"})
	vars, err := ScopeVars(es, "db")
	if err != nil {
		t.Fatalf("ScopeVars: %v", err)
	}
	if vars["DB_HOST"] != "localhost" || vars["DB_PORT"] != "5432" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestListScopes(t *testing.T) {
	es := baseScopeSet(t)
	_ = CreateScope(es, "db", []string{"DB_HOST"})
	_ = CreateScope(es, "api", []string{"API_KEY"})
	names := ListScopes(es)
	if len(names) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(names))
	}
}

func TestCreateScope_NilEnvSet(t *testing.T) {
	if err := CreateScope(nil, "db", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
