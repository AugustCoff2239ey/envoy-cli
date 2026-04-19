package envset

import (
	"testing"
)

func baseSqueezeSet() *EnvSet {
	es, _ := New("squeeze-test", "local")
	_ = es.Set("HOST", "localhost")
	_ = es.Set("DB_HOST", "localhost") // duplicate of HOST
	_ = es.Set("PORT", "5432")
	_ = es.Set("DB_PORT", "5432") // duplicate of PORT
	_ = es.Set("NAME", "myapp")   // unique
	return es
}

func TestSqueeze_KeepFirst(t *testing.T) {
	es := baseSqueezeSet()
	res, err := Squeeze(es, DefaultSqueezeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d: %v", len(res.Removed), res.Removed)
	}
	// DB_HOST and DB_PORT come after HOST and PORT alphabetically, so they are removed
	if _, ok := es.Vars["HOST"]; !ok {
		t.Error("HOST should be retained")
	}
	if _, ok := es.Vars["PORT"]; !ok {
		t.Error("PORT should be retained")
	}
	if _, ok := es.Vars["NAME"]; !ok {
		t.Error("NAME should be retained")
	}
}

func TestSqueeze_KeepLast(t *testing.T) {
	es := baseSqueezeSet()
	opts := SqueezeOptions{KeepFirst: false}
	res, err := Squeeze(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestSqueeze_UniqueValues(t *testing.T) {
	es, _ := New("sq", "local")
	_ = es.Set("A", "alpha")
	_ = es.Set("B", "beta")
	_ = es.Set("C", "gamma")
	res, err := Squeeze(es, DefaultSqueezeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
	if len(es.Vars) != 3 {
		t.Errorf("expected 3 vars, got %d", len(es.Vars))
	}
}

func TestSqueeze_NilEnvSet(t *testing.T) {
	_, err := Squeeze(nil, DefaultSqueezeOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestSqueeze_EmptySet(t *testing.T) {
	es, _ := New("empty", "local")
	res, err := Squeeze(es, DefaultSqueezeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals on empty set")
	}
}
