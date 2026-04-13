package envset

import (
	"testing"
)

func baseTagSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("myapp", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return es
}

func TestNewTag_Valid(t *testing.T) {
	tag, err := NewTag("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "production" {
		t.Errorf("expected 'production', got %q", tag.Name)
	}
}

func TestNewTag_Empty(t *testing.T) {
	_, err := NewTag("")
	if err == nil {
		t.Fatal("expected error for empty tag name")
	}
}

func TestNewTag_InvalidChars(t *testing.T) {
	_, err := NewTag("bad tag!")
	if err == nil {
		t.Fatal("expected error for tag with spaces/special chars")
	}
}

func TestAddTag_AndHasTag(t *testing.T) {
	es := baseTagSet(t)
	tag, _ := NewTag("infra")
	if err := AddTag(es, tag); err != nil {
		t.Fatalf("AddTag: %v", err)
	}
	if !HasTag(es, "infra") {
		t.Error("expected envset to have tag 'infra'")
	}
}

func TestAddTag_Duplicate(t *testing.T) {
	es := baseTagSet(t)
	tag, _ := NewTag("infra")
	_ = AddTag(es, tag)
	err := AddTag(es, tag)
	if err == nil {
		t.Fatal("expected error when adding duplicate tag")
	}
}

func TestAddTag_NilEnvSet(t *testing.T) {
	tag, _ := NewTag("infra")
	if err := AddTag(nil, tag); err == nil {
		t.Fatal("expected error for nil envset")
	}
}

func TestRemoveTag(t *testing.T) {
	es := baseTagSet(t)
	tag, _ := NewTag("infra")
	_ = AddTag(es, tag)
	if err := RemoveTag(es, "infra"); err != nil {
		t.Fatalf("RemoveTag: %v", err)
	}
	if HasTag(es, "infra") {
		t.Error("expected tag to be removed")
	}
}

func TestRemoveTag_NotFound(t *testing.T) {
	es := baseTagSet(t)
	if err := RemoveTag(es, "ghost"); err == nil {
		t.Fatal("expected error when removing non-existent tag")
	}
}

func TestTagsSorted(t *testing.T) {
	es := baseTagSet(t)
	for _, name := range []string{"zebra", "alpha", "middle"} {
		tag, _ := NewTag(name)
		_ = AddTag(es, tag)
	}
	expected := []string{"alpha", "middle", "zebra"}
	for i, got := range es.Tags {
		if got != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], got)
		}
	}
}
