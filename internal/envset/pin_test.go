package envset

import (
	"testing"
)

func basePinSet() *EnvSet {
	es, _ := New("pintest", "staging")
	_ = es.Set("API_KEY", "abc123")
	_ = es.Set("DB_URL", "postgres://localhost/db")
	_ = es.Set("DEBUG", "true")
	return es
}

func TestPinKey_Valid(t *testing.T) {
	es := basePinSet()
	entry, err := PinKey(es, "API_KEY", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", entry.Key)
	}
	if entry.Value != "abc123" {
		t.Errorf("expected value abc123, got %s", entry.Value)
	}
	if entry.PinnedBy != "alice" {
		t.Errorf("expected pinnedBy alice, got %s", entry.PinnedBy)
	}
	if !IsPinned(es, "API_KEY") {
		t.Error("expected API_KEY to be pinned")
	}
}

func TestPinKey_NonExistentKey(t *testing.T) {
	es := basePinSet()
	_, err := PinKey(es, "MISSING_KEY", "bob")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPinKey_InvalidKey(t *testing.T) {
	es := basePinSet()
	_, err := PinKey(es, "123INVALID", "bob")
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestPinKey_NilEnvSet(t *testing.T) {
	_, err := PinKey(nil, "API_KEY", "bob")
	if err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestUnpinKey_Valid(t *testing.T) {
	es := basePinSet()
	_, _ = PinKey(es, "DB_URL", "carol")
	if err := UnpinKey(es, "DB_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsPinned(es, "DB_URL") {
		t.Error("expected DB_URL to be unpinned")
	}
}

func TestUnpinKey_NotPinned(t *testing.T) {
	es := basePinSet()
	err := UnpinKey(es, "DEBUG")
	if err == nil {
		t.Fatal("expected error when unpinning a non-pinned key")
	}
}

func TestPinnedKeys_Multiple(t *testing.T) {
	es := basePinSet()
	_, _ = PinKey(es, "API_KEY", "alice")
	_, _ = PinKey(es, "DB_URL", "alice")
	pinned := PinnedKeys(es)
	if len(pinned) != 2 {
		t.Errorf("expected 2 pinned keys, got %d", len(pinned))
	}
}

func TestIsPinned_NilEnvSet(t *testing.T) {
	if IsPinned(nil, "API_KEY") {
		t.Error("expected false for nil EnvSet")
	}
}
