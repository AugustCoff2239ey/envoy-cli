package envset

import (
	"testing"
)

func basePatchSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("patch-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("APP_HOST", "localhost")
	_ = es.Set("APP_PORT", "8080")
	_ = es.Set("APP_DEBUG", "true")
	return es
}

func TestPatch_SetAndDelete(t *testing.T) {
	es := basePatchSet(t)
	entries := []PatchEntry{
		{Op: PatchOpSet, Key: "APP_PORT", Value: "9090"},
		{Op: PatchOpDelete, Key: "APP_DEBUG"},
		{Op: PatchOpSet, Key: "APP_NEW", Value: "added"},
	}
	if err := Patch(es, entries); err != nil {
		t.Fatalf("Patch: %v", err)
	}
	if v, _ := es.Get("APP_PORT"); v != "9090" {
		t.Errorf("expected APP_PORT=9090, got %q", v)
	}
	if _, ok := es.Get("APP_DEBUG"); ok {
		t.Error("expected APP_DEBUG to be deleted")
	}
	if v, _ := es.Get("APP_NEW"); v != "added" {
		t.Errorf("expected APP_NEW=added, got %q", v)
	}
}

func TestPatch_NilEnvSet(t *testing.T) {
	err := Patch(nil, []PatchEntry{{Op: PatchOpSet, Key: "X", Value: "1"}})
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestPatch_InvalidKey(t *testing.T) {
	es := basePatchSet(t)
	err := Patch(es, []PatchEntry{{Op: PatchOpSet, Key: "bad key!", Value: "v"}})
	if err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestPatch_LockedKey(t *testing.T) {
	es := basePatchSet(t)
	_ = LockKey(es, "APP_HOST")
	err := Patch(es, []PatchEntry{{Op: PatchOpSet, Key: "APP_HOST", Value: "remote"}})
	if err == nil {
		t.Error("expected error for locked key")
	}
}

func TestPatch_UnknownOp(t *testing.T) {
	es := basePatchSet(t)
	err := Patch(es, []PatchEntry{{Op: PatchOp("upsert"), Key: "APP_PORT", Value: "1"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatchFromDiff(t *testing.T) {
	base, _ := New("base", "local")
	_ = base.Set("KEEP", "same")
	_ = base.Set("CHANGE", "old")
	_ = base.Set("REMOVE", "gone")

	target, _ := New("target", "local")
	_ = target.Set("KEEP", "same")
	_ = target.Set("CHANGE", "new")
	_ = target.Set("ADD", "fresh")

	d := Diff(base, target)
	entries := PatchFromDiff(d)

	if err := Patch(base, entries); err != nil {
		t.Fatalf("Patch: %v", err)
	}
	if v, _ := base.Get("CHANGE"); v != "new" {
		t.Errorf("expected CHANGE=new, got %q", v)
	}
	if v, _ := base.Get("ADD"); v != "fresh" {
		t.Errorf("expected ADD=fresh, got %q", v)
	}
	if _, ok := base.Get("REMOVE"); ok {
		t.Error("expected REMOVE to be deleted")
	}
}
