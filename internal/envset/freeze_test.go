package envset

import (
	"testing"
)

func baseFreezeSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("freeze-test", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("API_KEY", "abc123")
	return es
}

func TestFreeze_Valid(t *testing.T) {
	es := baseFreezeSet(t)
	if err := Freeze(es); err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	if !IsFrozen(es) {
		t.Error("expected envset to be frozen")
	}
}

func TestUnfreeze_Valid(t *testing.T) {
	es := baseFreezeSet(t)
	_ = Freeze(es)
	if err := Unfreeze(es); err != nil {
		t.Fatalf("Unfreeze: %v", err)
	}
	if IsFrozen(es) {
		t.Error("expected envset to be unfrozen")
	}
}

func TestIsFrozen_Default(t *testing.T) {
	es := baseFreezeSet(t)
	if IsFrozen(es) {
		t.Error("new envset should not be frozen")
	}
}

func TestIsFrozen_NilEnvSet(t *testing.T) {
	if IsFrozen(nil) {
		t.Error("nil envset should not be frozen")
	}
}

func TestAssertMutable_WhenFrozen(t *testing.T) {
	es := baseFreezeSet(t)
	_ = Freeze(es)
	if err := AssertMutable(es); err == nil {
		t.Error("expected error for frozen envset")
	}
}

func TestAssertMutable_WhenNotFrozen(t *testing.T) {
	es := baseFreezeSet(t)
	if err := AssertMutable(es); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFreeze_NilEnvSet(t *testing.T) {
	if err := Freeze(nil); err == nil {
		t.Error("expected error for nil envset")
	}
}
