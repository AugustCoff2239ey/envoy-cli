package envset

import (
	"strings"
	"testing"
)

func baseSplitSet(t *testing.T) *EnvSet {
	t.Helper()
	es, _ := New("app", "staging")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("APP_PORT", "8080")
	_ = es.Set("APP_DEBUG", "true")
	_ = es.Set("LOG_LEVEL", "info")
	return es
}

func TestSplit_ByKeys(t *testing.T) {
	src := baseSplitSet(t)
	res, err := Split(src, SplitOptions{Keys: []string{"DB_HOST", "DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Matched.Vars["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in matched")
	}
	if _, ok := res.Unmatched.Vars["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in unmatched")
	}
	if len(res.Matched.Vars)+len(res.Unmatched.Vars) != len(src.Vars) {
		t.Error("total key count mismatch")
	}
}

func TestSplit_ByPredicate(t *testing.T) {
	src := baseSplitSet(t)
	res, err := Split(src, SplitOptions{
		Predicate: func(k, _ string) bool { return strings.HasPrefix(k, "DB_") },
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Matched.Vars) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched.Vars))
	}
	if len(res.Unmatched.Vars) != 3 {
		t.Errorf("expected 3 unmatched, got %d", len(res.Unmatched.Vars))
	}
}

func TestSplit_CustomNames(t *testing.T) {
	src := baseSplitSet(t)
	res, err := Split(src, SplitOptions{
		Keys:          []string{"LOG_LEVEL"},
		MatchedName:   "logging",
		UnmatchedName: "rest",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched.Name != "logging" {
		t.Errorf("expected matched name 'logging', got %q", res.Matched.Name)
	}
	if res.Unmatched.Name != "rest" {
		t.Errorf("expected unmatched name 'rest', got %q", res.Unmatched.Name)
	}
}

func TestSplit_NilSource(t *testing.T) {
	_, err := Split(nil, SplitOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestSplit_NoPredicate_AllUnmatched(t *testing.T) {
	src := baseSplitSet(t)
	res, err := Split(src, SplitOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Matched.Vars) != 0 {
		t.Errorf("expected 0 matched, got %d", len(res.Matched.Vars))
	}
	if len(res.Unmatched.Vars) != len(src.Vars) {
		t.Errorf("expected all keys unmatched")
	}
}
