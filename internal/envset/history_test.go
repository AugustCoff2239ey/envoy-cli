package envset

import (
	"testing"
)

func TestNewHistory_Empty(t *testing.T) {
	h := NewHistory()
	if h == nil {
		t.Fatal("expected non-nil History")
	}
	if len(h.Entries()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(h.Entries()))
	}
}

func TestHistory_Record_Valid(t *testing.T) {
	h := NewHistory()
	err := h.Record("set", "DB_HOST", "", "localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != "set" || e.Key != "DB_HOST" || e.NewValue != "localhost" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestHistory_Record_EmptyAction(t *testing.T) {
	h := NewHistory()
	err := h.Record("", "KEY", "", "val")
	if err == nil {
		t.Error("expected error for empty action")
	}
}

func TestHistory_Record_NilReceiver(t *testing.T) {
	var h *History
	err := h.Record("set", "KEY", "", "val")
	if err == nil {
		t.Error("expected error for nil receiver")
	}
}

func TestHistory_FilterByKey(t *testing.T) {
	h := NewHistory()
	_ = h.Record("set", "DB_HOST", "", "localhost")
	_ = h.Record("set", "API_KEY", "", "abc123")
	_ = h.Record("delete", "DB_HOST", "localhost", "")

	results := h.FilterByKey("DB_HOST")
	if len(results) != 2 {
		t.Errorf("expected 2 entries for DB_HOST, got %d", len(results))
	}
}

func TestHistory_FilterByAction(t *testing.T) {
	h := NewHistory()
	_ = h.Record("set", "DB_HOST", "", "localhost")
	_ = h.Record("set", "API_KEY", "", "secret")
	_ = h.Record("delete", "DB_HOST", "localhost", "")

	results := h.FilterByAction("set")
	if len(results) != 2 {
		t.Errorf("expected 2 'set' entries, got %d", len(results))
	}
}

func TestHistory_Clear(t *testing.T) {
	h := NewHistory()
	_ = h.Record("set", "KEY", "", "value")
	h.Clear()
	if len(h.Entries()) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(h.Entries()))
	}
}

func TestHistory_Entries_ReturnsCopy(t *testing.T) {
	h := NewHistory()
	_ = h.Record("set", "X", "", "1")
	copy1 := h.Entries()
	copy1[0].Key = "MUTATED"
	copy2 := h.Entries()
	if copy2[0].Key == "MUTATED" {
		t.Error("Entries() should return a copy, not a reference")
	}
}
