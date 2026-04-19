package envset

import (
	"testing"
)

func baseReorderSet() *EnvSet {
	es, _ := New("reorder-test", "local")
	_ = es.Set("ALPHA", "1")
	_ = es.Set("BETA", "2")
	_ = es.Set("GAMMA", "3")
	return es
}

func TestReorder_ExplicitOrder(t *testing.T) {
	es := baseReorderSet()
	res, err := Reorder(es, []string{"GAMMA", "ALPHA", "BETA"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Keys[0] != "GAMMA" || res.Keys[1] != "ALPHA" || res.Keys[2] != "BETA" {
		t.Errorf("unexpected order: %v", res.Keys)
	}
}

func TestReorder_PartialOrder(t *testing.T) {
	es := baseReorderSet()
	res, err := Reorder(es, []string{"GAMMA"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Keys[0] != "GAMMA" {
		t.Errorf("expected GAMMA first, got %v", res.Keys[0])
	}
	if len(res.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Keys))
	}
}

func TestReorder_MissingKey(t *testing.T) {
	es := baseReorderSet()
	_, err := Reorder(es, []string{"MISSING"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestReorder_NilEnvSet(t *testing.T) {
	_, err := Reorder(nil, []string{"ALPHA"})
	if err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestReorder_EmptyOrder(t *testing.T) {
	es := baseReorderSet()
	_, err := Reorder(es, []string{})
	if err == nil {
		t.Fatal("expected error for empty order")
	}
}

func TestGetOrder_AfterReorder(t *testing.T) {
	es := baseReorderSet()
	_, _ = Reorder(es, []string{"BETA", "ALPHA", "GAMMA"})
	order := GetOrder(es)
	if len(order) != 3 || order[0] != "BETA" {
		t.Errorf("unexpected stored order: %v", order)
	}
}

func TestGetOrder_NoOrder(t *testing.T) {
	es := baseReorderSet()
	if order := GetOrder(es); order != nil {
		t.Errorf("expected nil order, got %v", order)
	}
}
