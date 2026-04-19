package envset

import (
	"testing"
	"time"
)

func baseStatusSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("status-set", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_URL", "postgres://localhost")
	_ = es.Set("API_KEY", "secret")
	_ = es.Set("PORT", "8080")
	return es
}

func TestStatus_Basic(t *testing.T) {
	es := baseStatusSet(t)
	r, err := Status(es)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 keys, got %d", r.TotalKeys)
	}
	if r.Name != "status-set" || r.Environment != "staging" {
		t.Errorf("unexpected name/env: %s/%s", r.Name, r.Environment)
	}
	if r.Readonly || r.Frozen {
		t.Error("expected not readonly or frozen")
	}
}

func TestStatus_LockedAndPinned(t *testing.T) {
	es := baseStatusSet(t)
	_, _ = LockKey(es, "DB_URL", "test")
	_ = PinKey(es, "API_KEY")
	r, err := Status(es)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if r.LockedKeys != 1 {
		t.Errorf("expected 1 locked, got %d", r.LockedKeys)
	}
	if r.PinnedKeys != 1 {
		t.Errorf("expected 1 pinned, got %d", r.PinnedKeys)
	}
}

func TestStatus_Expired(t *testing.T) {
	es := baseStatusSet(t)
	past := time.Now().Add(-time.Hour)
	_ = SetExpiry(es, "PORT", past)
	r, err := Status(es)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if r.ExpiredKeys != 1 {
		t.Errorf("expected 1 expired, got %d", r.ExpiredKeys)
	}
}

func TestStatus_ReadonlyFrozen(t *testing.T) {
	es := baseStatusSet(t)
	MarkReadonly(es)
	Freeze(es)
	r, err := Status(es)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !r.Readonly {
		t.Error("expected readonly")
	}
	if !r.Frozen {
		t.Error("expected frozen")
	}
}

func TestStatus_NilEnvSet(t *testing.T) {
	_, err := Status(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
