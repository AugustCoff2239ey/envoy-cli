package envset

import (
	"testing"
)

func baseInheritSets(t *testing.T) (parent, child *EnvSet) {
	t.Helper()
	parent, _ = New("parent", "production")
	parent.Vars["DB_HOST"] = "prod-db.example.com"
	parent.Vars["DB_PORT"] = "5432"
	parent.Vars["API_KEY"] = "secret-prod"

	child, _ = New("child", "staging")
	child.Vars["DB_HOST"] = "staging-db.example.com" // already set
	return parent, child
}

func TestInherit_AllKeys(t *testing.T) {
	parent, child := baseInheritSets(t)
	res, err := Inherit(parent, child, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Inherited) != 2 {
		t.Errorf("expected 2 inherited, got %d", len(res.Inherited))
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST skipped, got %v", res.Skipped)
	}
	if child.Vars["DB_HOST"] != "staging-db.example.com" {
		t.Error("child DB_HOST should not be overwritten")
	}
	if child.Vars["DB_PORT"] != "5432" {
		t.Error("expected DB_PORT to be inherited")
	}
}

func TestInherit_SelectedKeys(t *testing.T) {
	parent, child := baseInheritSets(t)
	res, err := Inherit(parent, child, []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Inherited) != 1 || res.Inherited[0] != "API_KEY" {
		t.Errorf("expected API_KEY inherited, got %v", res.Inherited)
	}
	if _, ok := child.Vars["DB_PORT"]; ok {
		t.Error("DB_PORT should not have been inherited")
	}
}

func TestInherit_MissingParentKey(t *testing.T) {
	parent, child := baseInheritSets(t)
	_, err := Inherit(parent, child, []string{"NONEXISTENT"})
	if err == nil {
		t.Fatal("expected error for missing parent key")
	}
}

func TestInherit_NilParent(t *testing.T) {
	_, child := baseInheritSets(t)
	_, err := Inherit(nil, child, nil)
	if err == nil {
		t.Fatal("expected error for nil parent")
	}
}

func TestInherit_NilChild(t *testing.T) {
	parent, _ := baseInheritSets(t)
	_, err := Inherit(parent, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil child")
	}
}

func TestInherit_InvalidKey(t *testing.T) {
	parent, _ := New("parent", "production")
	parent.Vars["bad key"] = "value"
	child, _ := New("child", "staging")
	_, err := Inherit(parent, child, []string{"bad key"})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}
