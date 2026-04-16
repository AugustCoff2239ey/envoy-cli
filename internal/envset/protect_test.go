package envset

import (
	"testing"
)

func baseProtectSet(t *testing.T) *EnvSet {
	t.Helper()
	e, err := New("protect-test", "local")
	if err != nil {
		t.Fatal(err)
	}
	_ = e.Set("API_KEY", "secret")
	_ = e.Set("DB_PASS", "pass123")
	return e
}

func TestProtectKey_Valid(t *testing.T) {
	e := baseProtectSet(t)
	if err := ProtectKey(e, "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsProtected(e, "API_KEY") {
		t.Error("expected API_KEY to be protected")
	}
}

func TestProtectKey_AlreadyProtected(t *testing.T) {
	e := baseProtectSet(t)
	_ = ProtectKey(e, "API_KEY")
	if err := ProtectKey(e, "API_KEY"); err != ErrAlreadyProtected {
		t.Errorf("expected ErrAlreadyProtected, got %v", err)
	}
}

func TestProtectKey_NonExistentKey(t *testing.T) {
	e := baseProtectSet(t)
	if err := ProtectKey(e, "MISSING"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestProtectKey_NilEnvSet(t *testing.T) {
	if err := ProtectKey(nil, "API_KEY"); err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestUnprotectKey_Valid(t *testing.T) {
	e := baseProtectSet(t)
	_ = ProtectKey(e, "API_KEY")
	if err := UnprotectKey(e, "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsProtected(e, "API_KEY") {
		t.Error("expected API_KEY to be unprotected")
	}
}

func TestUnprotectKey_NotProtected(t *testing.T) {
	e := baseProtectSet(t)
	if err := UnprotectKey(e, "API_KEY"); err != ErrNotProtected {
		t.Errorf("expected ErrNotProtected, got %v", err)
	}
}

func TestProtectedKeys_List(t *testing.T) {
	e := baseProtectSet(t)
	_ = ProtectKey(e, "API_KEY")
	_ = ProtectKey(e, "DB_PASS")
	keys := ProtectedKeys(e)
	if len(keys) != 2 {
		t.Errorf("expected 2 protected keys, got %d", len(keys))
	}
}
