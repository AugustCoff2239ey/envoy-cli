package envset

import (
	"testing"
)

func baseLockSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("locktest", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_URL", "postgres://localhost/db")
	_ = es.Set("API_KEY", "secret")
	_ = es.Set("PORT", "8080")
	return es
}

func TestLockKey_Valid(t *testing.T) {
	es := baseLockSet(t)
	if err := LockKey(es, "DB_URL", "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsLocked(es, "DB_URL") {
		t.Error("expected DB_URL to be locked")
	}
}

func TestLockKey_NonExistentKey(t *testing.T) {
	es := baseLockSet(t)
	if err := LockKey(es, "MISSING_KEY", "admin"); err == nil {
		t.Error("expected error for non-existent key")
	}
}

func TestLockKey_InvalidKey(t *testing.T) {
	es := baseLockSet(t)
	if err := LockKey(es, "invalid key!", "admin"); err == nil {
		t.Error("expected error for invalid key name")
	}
}

func TestLockKey_NilEnvSet(t *testing.T) {
	if err := LockKey(nil, "DB_URL", "admin"); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestUnlockKey_Valid(t *testing.T) {
	es := baseLockSet(t)
	_ = LockKey(es, "API_KEY", "ci")
	if err := UnlockKey(es, "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsLocked(es, "API_KEY") {
		t.Error("expected API_KEY to be unlocked")
	}
}

func TestUnlockKey_NotLocked(t *testing.T) {
	es := baseLockSet(t)
	if err := UnlockKey(es, "PORT"); err == nil {
		t.Error("expected error when unlocking a non-locked key")
	}
}

func TestIsLocked_NilEnvSet(t *testing.T) {
	if IsLocked(nil, "KEY") {
		t.Error("expected false for nil envset")
	}
}

func TestLockedKeys_ReturnsAll(t *testing.T) {
	es := baseLockSet(t)
	_ = LockKey(es, "DB_URL", "admin")
	_ = LockKey(es, "API_KEY", "ci")
	entries := LockedKeys(es)
	if len(entries) != 2 {
		t.Errorf("expected 2 locked keys, got %d", len(entries))
	}
}

func TestLockedKeys_NilEnvSet(t *testing.T) {
	if entries := LockedKeys(nil); entries != nil {
		t.Error("expected nil for nil envset")
	}
}
