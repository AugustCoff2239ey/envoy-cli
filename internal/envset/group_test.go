package envset

import (
	"testing"
)

func baseGroupSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("group-test", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	_ = es.Set("LOG_LEVEL", "info")
	return es
}

func TestCreateGroup_Valid(t *testing.T) {
	es := baseGroupSet(t)
	if err := CreateGroup(es, "database", []string{"DB_HOST", "DB_PORT"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, err := GetGroup(es, "database")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(g.Keys))
	}
}

func TestCreateGroup_EmptyName(t *testing.T) {
	es := baseGroupSet(t)
	if err := CreateGroup(es, "", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for empty group name")
	}
}

func TestCreateGroup_InvalidName(t *testing.T) {
	es := baseGroupSet(t)
	if err := CreateGroup(es, "my group!", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for invalid group name")
	}
}

func TestCreateGroup_MissingKey(t *testing.T) {
	es := baseGroupSet(t)
	if err := CreateGroup(es, "db", []string{"DB_HOST", "MISSING_KEY"}); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestCreateGroup_NilEnvSet(t *testing.T) {
	if err := CreateGroup(nil, "db", []string{"DB_HOST"}); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestListGroups(t *testing.T) {
	es := baseGroupSet(t)
	_ = CreateGroup(es, "database", []string{"DB_HOST", "DB_PORT"})
	_ = CreateGroup(es, "app", []string{"API_KEY", "LOG_LEVEL"})
	groups := ListGroups(es)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "app" || groups[1].Name != "database" {
		t.Errorf("groups not sorted: %v", groups)
	}
}

func TestDeleteGroup_Valid(t *testing.T) {
	es := baseGroupSet(t)
	_ = CreateGroup(es, "database", []string{"DB_HOST"})
	if err := DeleteGroup(es, "database"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := GetGroup(es, "database"); err == nil {
		t.Error("expected group to be deleted")
	}
}

func TestDeleteGroup_NotFound(t *testing.T) {
	es := baseGroupSet(t)
	if err := DeleteGroup(es, "nonexistent"); err == nil {
		t.Error("expected error for nonexistent group")
	}
}

func TestListGroups_NilEnvSet(t *testing.T) {
	if groups := ListGroups(nil); groups != nil {
		t.Error("expected nil for nil envset")
	}
}
