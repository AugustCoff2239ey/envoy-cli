package envset

import (
	"testing"
)

func baseNoteSet() *EnvSet {
	es, _ := New("notes-test", "staging")
	_ = es.Set("API_KEY", "abc123")
	return es
}

func TestAddNote_Valid(t *testing.T) {
	es := baseNoteSet()
	if err := AddNote(es, "Initial setup", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	notes := ListNotes(es)
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Text != "Initial setup" {
		t.Errorf("expected 'Initial setup', got %q", notes[0].Text)
	}
	if notes[0].Author != "alice" {
		t.Errorf("expected author 'alice', got %q", notes[0].Author)
	}
}

func TestAddNote_MultipleNotes(t *testing.T) {
	es := baseNoteSet()
	_ = AddNote(es, "First note", "alice")
	_ = AddNote(es, "Second note", "bob")
	notes := ListNotes(es)
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
	if notes[1].Text != "Second note" {
		t.Errorf("expected 'Second note', got %q", notes[1].Text)
	}
}

func TestAddNote_EmptyText(t *testing.T) {
	es := baseNoteSet()
	if err := AddNote(es, "", "alice"); err == nil {
		t.Error("expected error for empty text")
	}
}

func TestAddNote_NewlineInText(t *testing.T) {
	es := baseNoteSet()
	if err := AddNote(es, "bad\nnote", "alice"); err == nil {
		t.Error("expected error for newline in text")
	}
}

func TestAddNote_NilEnvSet(t *testing.T) {
	if err := AddNote(nil, "hello", "alice"); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestClearNotes_Valid(t *testing.T) {
	es := baseNoteSet()
	_ = AddNote(es, "Note one", "alice")
	_ = AddNote(es, "Note two", "bob")
	if err := ClearNotes(es); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notes := ListNotes(es); len(notes) != 0 {
		t.Errorf("expected 0 notes after clear, got %d", len(notes))
	}
}

func TestListNotes_NilEnvSet(t *testing.T) {
	if notes := ListNotes(nil); notes != nil {
		t.Errorf("expected nil for nil envset, got %v", notes)
	}
}
