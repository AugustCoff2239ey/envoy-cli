package envset

import (
	"testing"
)

func baseAliasSet() *EnvSet {
	es, _ := New("alias-test", "local")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	return es
}

func TestAddAlias_Valid(t *testing.T) {
	es := baseAliasSet()
	if err := AddAlias(es, "DB_HOST", "database-host"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aliases := ListAliases(es, "DB_HOST")
	if len(aliases) != 1 || aliases[0] != "database-host" {
		t.Errorf("expected alias 'database-host', got %v", aliases)
	}
}

func TestAddAlias_InvalidName(t *testing.T) {
	es := baseAliasSet()
	if err := AddAlias(es, "DB_HOST", "123bad!"); err == nil {
		t.Error("expected error for invalid alias name")
	}
}

func TestAddAlias_NonExistentKey(t *testing.T) {
	es := baseAliasSet()
	if err := AddAlias(es, "MISSING", "m"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestAddAlias_NilEnvSet(t *testing.T) {
	if err := AddAlias(nil, "DB_HOST", "host"); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestRemoveAlias_Valid(t *testing.T) {
	es := baseAliasSet()
	_ = AddAlias(es, "DB_HOST", "host")
	if err := RemoveAlias(es, "DB_HOST", "host"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ListAliases(es, "DB_HOST")) != 0 {
		t.Error("expected no aliases after removal")
	}
}

func TestRemoveAlias_NotFound(t *testing.T) {
	es := baseAliasSet()
	if err := RemoveAlias(es, "DB_HOST", "ghost"); err == nil {
		t.Error("expected error for non-existent alias")
	}
}

func TestResolveAlias_Found(t *testing.T) {
	es := baseAliasSet()
	_ = AddAlias(es, "DB_PORT", "port")
	key, ok := ResolveAlias(es, "port")
	if !ok || key != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %q ok=%v", key, ok)
	}
}

func TestResolveAlias_NotFound(t *testing.T) {
	es := baseAliasSet()
	_, ok := ResolveAlias(es, "nonexistent")
	if ok {
		t.Error("expected not found")
	}
}
