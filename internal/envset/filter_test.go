package envset

import (
	"testing"
)

func baseFilterSet(t *testing.T) *EnvSet {
	t.Helper()
	es, _ := New("filter-test", "staging")
	_ = es.Set("APP_HOST", "localhost")
	_ = es.Set("APP_PORT", "8080")
	_ = es.Set("DB_HOST", "db.local")
	_ = es.Set("DB_PASSWORD", "secret")
	_ = es.Set("LOG_LEVEL", "debug")
	return es
}

func TestFilter_ByPrefix(t *testing.T) {
	es := baseFilterSet(t)
	out, err := Filter(es, FilterOptions{Prefix: "APP"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Vars["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := out.Vars["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST")
	}
}

func TestFilter_BySuffix(t *testing.T) {
	es := baseFilterSet(t)
	out, err := Filter(es, FilterOptions{Suffix: "HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Vars) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out.Vars))
	}
}

func TestFilter_Exclude(t *testing.T) {
	es := baseFilterSet(t)
	out, err := Filter(es, FilterOptions{Exclude: []string{"DB_PASSWORD", "LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Vars["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if len(out.Vars) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out.Vars))
	}
}

func TestFilter_EnvMismatch(t *testing.T) {
	es := baseFilterSet(t)
	out, err := Filter(es, FilterOptions{Envs: []string{"production"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Vars) != 0 {
		t.Errorf("expected empty result for env mismatch, got %d keys", len(out.Vars))
	}
}

func TestFilter_NilEnvSet(t *testing.T) {
	_, err := Filter(nil, FilterOptions{})
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
