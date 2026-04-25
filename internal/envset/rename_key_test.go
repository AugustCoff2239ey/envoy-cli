package envset

import (
	"testing"
)

func baseRenameKeySet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("rename-key-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["OLD_KEY"] = "hello"
	es.Vars["OTHER_KEY"] = "world"
	return es
}

func TestRenameKey_Valid(t *testing.T) {
	es := baseRenameKeySet(t)
	if err := RenameKey(es, "OLD_KEY", "NEW_KEY", DefaultRenameKeyOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := es.Vars["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if v, ok := es.Vars["NEW_KEY"]; !ok || v != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", v)
	}
}

func TestRenameKey_TargetExists_NoOverwrite(t *testing.T) {
	es := baseRenameKeySet(t)
	if err := RenameKey(es, "OLD_KEY", "OTHER_KEY", DefaultRenameKeyOptions()); err == nil {
		t.Fatal("expected error when target exists without overwrite")
	}
}

func TestRenameKey_TargetExists_WithOverwrite(t *testing.T) {
	es := baseRenameKeySet(t)
	opts := RenameKeyOptions{Overwrite: true}
	if err := RenameKey(es, "OLD_KEY", "OTHER_KEY", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := es.Vars["OTHER_KEY"]; v != "hello" {
		t.Errorf("expected OTHER_KEY=hello after overwrite, got %q", v)
	}
}

func TestRenameKey_NotFound(t *testing.T) {
	es := baseRenameKeySet(t)
	if err := RenameKey(es, "MISSING", "NEW_KEY", DefaultRenameKeyOptions()); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKey_LockedKey(t *testing.T) {
	es := baseRenameKeySet(t)
	if err := LockKey(es, "OLD_KEY", "reason"); err != nil {
		t.Fatalf("LockKey: %v", err)
	}
	if err := RenameKey(es, "OLD_KEY", "NEW_KEY", DefaultRenameKeyOptions()); err == nil {
		t.Fatal("expected error renaming locked key")
	}
}

func TestRenameKey_NilEnvSet(t *testing.T) {
	if err := RenameKey(nil, "OLD_KEY", "NEW_KEY", DefaultRenameKeyOptions()); err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestRenameKey_InvalidNewKey(t *testing.T) {
	es := baseRenameKeySet(t)
	if err := RenameKey(es, "OLD_KEY", "123INVALID", DefaultRenameKeyOptions()); err == nil {
		t.Fatal("expected error for invalid new key name")
	}
}
