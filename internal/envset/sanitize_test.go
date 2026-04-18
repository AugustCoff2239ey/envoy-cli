package envset

import (
	"testing"
)

func baseSanitizeSet(t *testing.T) *EnvSet {
	t.Helper()
	e, _ := New("sanitize-test", "local")
	e.Vars["KEY_A"] = "  hello  "
	e.Vars["KEY_B"] = "world\x01\x02"
	e.Vars["KEY WITH SPACE"] = "value"
	return e
}

func TestSanitize_TrimWhitespace(t *testing.T) {
	e := baseSanitizeSet(t)
	opts := DefaultSanitizeOptions()
	opts.ReplaceSpacesInKeys = false

	_, err := Sanitize(e, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := e.Vars["KEY_A"]; got != "hello" {
		t.Errorf("expected trimmed value, got %q", got)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	e := baseSanitizeSet(t)
	opts := DefaultSanitizeOptions()

	_, err := Sanitize(e, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := e.Vars["KEY_B"]; got != "world" {
		t.Errorf("expected control chars stripped, got %q", got)
	}
}

func TestSanitize_ReplaceSpacesInKeys(t *testing.T) {
	e := baseSanitizeSet(t)
	opts := DefaultSanitizeOptions()

	_, err := Sanitize(e, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := e.Vars["KEY_WITH_SPACE"]; !ok {
		t.Error("expected space in key replaced with underscore")
	}
	if _, ok := e.Vars["KEY WITH SPACE"]; ok {
		t.Error("old key with space should be removed")
	}
}

func TestSanitize_NilEnvSet(t *testing.T) {
	_, err := Sanitize(nil, DefaultSanitizeOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestSanitize_ChangedFlag(t *testing.T) {
	e, _ := New("flag-test", "local")
	e.Vars["CLEAN_KEY"] = "clean_value"

	results, err := Sanitize(e, DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "CLEAN_KEY" && r.Changed {
			t.Error("clean key/value should not be marked as changed")
		}
	}
}
