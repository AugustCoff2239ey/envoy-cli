package envset

import (
	"testing"
)

func baseLabelSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("labeltest", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return es
}

func TestAddLabel_Valid(t *testing.T) {
	es := baseLabelSet(t)
	if err := AddLabel(es, "team", "platform"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetLabel(es, "team")
	if !ok || v != "platform" {
		t.Errorf("expected label team=platform, got %q %v", v, ok)
	}
}

func TestAddLabel_InvalidKey(t *testing.T) {
	es := baseLabelSet(t)
	if err := AddLabel(es, "bad key!", "v"); err == nil {
		t.Error("expected error for invalid label key")
	}
}

func TestAddLabel_NilEnvSet(t *testing.T) {
	if err := AddLabel(nil, "k", "v"); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestRemoveLabel_Valid(t *testing.T) {
	es := baseLabelSet(t)
	_ = AddLabel(es, "env", "prod")
	if err := RemoveLabel(es, "env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := GetLabel(es, "env"); ok {
		t.Error("label should have been removed")
	}
}

func TestRemoveLabel_NotFound(t *testing.T) {
	es := baseLabelSet(t)
	if err := RemoveLabel(es, "missing"); err == nil {
		t.Error("expected error for missing label")
	}
}

func TestListLabels(t *testing.T) {
	es := baseLabelSet(t)
	_ = AddLabel(es, "owner", "alice")
	_ = AddLabel(es, "tier", "free")
	labels := ListLabels(es)
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
	if labels["owner"] != "alice" || labels["tier"] != "free" {
		t.Errorf("unexpected labels: %v", labels)
	}
}

func TestListLabels_NilEnvSet(t *testing.T) {
	labels := ListLabels(nil)
	if len(labels) != 0 {
		t.Error("expected empty map for nil envset")
	}
}
