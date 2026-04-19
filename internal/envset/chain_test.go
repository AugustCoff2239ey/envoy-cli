package envset

import (
	"errors"
	"testing"
)

func baseChainSets(t *testing.T) []*EnvSet {
	t.Helper()
	a, _ := New("alpha", "local")
	_ = a.Set("FOO", "1")
	b, _ := New("beta", "staging")
	_ = b.Set("BAR", "2")
	c, _ := New("gamma", "production")
	_ = c.Set("BAZ", "3")
	return []*EnvSet{a, b, c}
}

func TestChain_AllApplied(t *testing.T) {
	sets := baseChainSets(t)
	result, err := Chain(sets, func(es *EnvSet) error {
		return es.Set("CHAINED", "yes")
	}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 3 {
		t.Errorf("expected 3 applied, got %d", len(result.Applied))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
	for _, es := range sets {
		if v, _ := es.Get("CHAINED"); v != "yes" {
			t.Errorf("%s: expected CHAINED=yes", es.Name)
		}
	}
}

func TestChain_StopOnError(t *testing.T) {
	sets := baseChainSets(t)
	called := 0
	_, err := Chain(sets, func(es *EnvSet) error {
		called++
		if es.Name == "beta" {
			return errors.New("forced error")
		}
		return nil
	}, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if called != 2 {
		t.Errorf("expected 2 calls before stop, got %d", called)
	}
}

func TestChain_ContinueOnError(t *testing.T) {
	sets := baseChainSets(t)
	result, err := Chain(sets, func(es *EnvSet) error {
		if es.Name == "beta" {
			return errors.New("forced error")
		}
		return nil
	}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(result.Applied))
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "beta" {
		t.Errorf("expected beta skipped, got %v", result.Skipped)
	}
}

func TestChain_NilFunction(t *testing.T) {
	sets := baseChainSets(t)
	_, err := Chain(sets, nil, false)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}

func TestChain_NilEnvSetSkipped(t *testing.T) {
	sets := []*EnvSet{nil, nil}
	result, err := Chain(sets, func(es *EnvSet) error { return nil }, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(result.Skipped))
	}
}
