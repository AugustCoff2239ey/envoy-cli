package envset

import (
	"testing"
)

func baseMergeSets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	dst, _ := New("dst", "local")
	_ = dst.Set("APP_NAME", "envoy")
	_ = dst.Set("LOG_LEVEL", "debug")
	_ = dst.Set("PORT", "8080")

	src, _ := New("src", "local")
	_ = src.Set("LOG_LEVEL", "info") // conflict
	_ = src.Set("DB_URL", "postgres://localhost/db") // new key
	return dst, src
}

func TestMerge_StrategyOurs(t *testing.T) {
	dst, src := baseMergeSets(t)
	res, err := Merge(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := res.Merged.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("expected 'debug', got %q", v)
	}
	if v, _ := res.Merged.Get("DB_URL"); v != "postgres://localhost/db" {
		t.Errorf("expected DB_URL to be merged in")
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "LOG_LEVEL" {
		t.Errorf("expected 1 conflict on LOG_LEVEL, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyTheirs(t *testing.T) {
	dst, src := baseMergeSets(t)
	res, err := Merge(dst, src, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := res.Merged.Get("LOG_LEVEL"); v != "info" {
		t.Errorf("expected 'info', got %q", v)
	}
}

func TestMerge_StrategyError(t *testing.T) {
	dst, src := baseMergeSets(t)
	_, err := Merge(dst, src, MergeStrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_NoConflict(t *testing.T) {
	dst, _ := New("dst", "local")
	_ = dst.Set("A", "1")
	src, _ := New("src", "local")
	_ = src.Set("B", "2")
	res, err := Merge(dst, src, MergeStrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if v, _ := res.Merged.Get("A"); v != "1" {
		t.Errorf("expected A=1")
	}
	if v, _ := res.Merged.Get("B"); v != "2" {
		t.Errorf("expected B=2")
	}
}

func TestMerge_NilInputs(t *testing.T) {
	_, err := Merge(nil, nil, MergeStrategyOurs)
	if err == nil {
		t.Fatal("expected error for nil inputs")
	}
}

func TestMerge_StrategyTheirs_PreservesNonConflicting(t *testing.T) {
	dst, src := baseMergeSets(t)
	res, err := Merge(dst, src, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Keys only in dst should still be present after a theirs-strategy merge.
	if v, _ := res.Merged.Get("APP_NAME"); v != "envoy" {
		t.Errorf("expected APP_NAME='envoy', got %q", v)
	}
	if v, _ := res.Merged.Get("PORT"); v != "8080" {
		t.Errorf("expected PORT='8080', got %q", v)
	}
	// New key from src should also be present.
	if v, _ := res.Merged.Get("DB_URL"); v != "postgres://localhost/db" {
		t.Errorf("expected DB_URL='postgres://localhost/db', got %q", v)
	}
}
