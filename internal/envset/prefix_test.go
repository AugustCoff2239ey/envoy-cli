package envset

import (
	"testing"
)

func basePrefixSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("prefix-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	return es
}

func TestAddPrefix_AllKeys(t *testing.T) {
	es := basePrefixSet(t)
	n, err := AddPrefix(es, "APP_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 renamed, got %d", n)
	}
	if _, ok := es.Vars["APP_DB_HOST"]; !ok {
		t.Error("expected APP_DB_HOST to exist")
	}
	if _, ok := es.Vars["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
}

func TestAddPrefix_SelectedKeys(t *testing.T) {
	es := basePrefixSet(t)
	n, err := AddPrefix(es, "X_", []string{"DB_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 renamed, got %d", n)
	}
	if _, ok := es.Vars["X_DB_HOST"]; !ok {
		t.Error("expected X_DB_HOST")
	}
	if _, ok := es.Vars["DB_PORT"]; !ok {
		t.Error("DB_PORT should be untouched")
	}
}

func TestAddPrefix_EmptyPrefix(t *testing.T) {
	es := basePrefixSet(t)
	_, err := AddPrefix(es, "", nil)
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestAddPrefix_NilEnvSet(t *testing.T) {
	_, err := AddPrefix(nil, "X_", nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestStripPrefix_AllKeys(t *testing.T) {
	es := basePrefixSet(t)
	n, err := StripPrefix(es, "DB_", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 stripped, got %d", n)
	}
	if _, ok := es.Vars["HOST"]; !ok {
		t.Error("expected HOST")
	}
	if _, ok := es.Vars["PORT"]; !ok {
		t.Error("expected PORT")
	}
	if _, ok := es.Vars["API_KEY"]; !ok {
		t.Error("API_KEY should remain")
	}
}

func TestStripPrefix_NilEnvSet(t *testing.T) {
	_, err := StripPrefix(nil, "DB_", nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
