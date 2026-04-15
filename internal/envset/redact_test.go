package envset

import (
	"testing"
)

func baseRedactSet() *EnvSet {
	es, _ := New("app", "production")
	_ = es.Set("DATABASE_URL", "postgres://user:secret@host/db")
	_ = es.Set("API_KEY", "abcdef123456")
	_ = es.Set("DEBUG", "false")
	_ = es.Set("EMPTY_VAR", "")
	return es
}

func TestRedact_FullMask(t *testing.T) {
	es := baseRedactSet()
	res, err := Redact(es, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, entry := range res.Entries {
		if entry.Key == "EMPTY_VAR" {
			if entry.Redacted != "[empty]" {
				t.Errorf("expected [empty] for empty value, got %q", entry.Redacted)
			}
			continue
		}
		for _, ch := range entry.Redacted {
			if ch != '*' {
				t.Errorf("key %s: expected all asterisks, got %q", entry.Key, entry.Redacted)
				break
			}
		}
	}
}

func TestRedact_RevealSuffix(t *testing.T) {
	es := baseRedactSet()
	res, err := Redact(es, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, entry := range res.Entries {
		if entry.Key == "API_KEY" {
			if entry.Redacted != "********3456" {
				t.Errorf("expected ********3456, got %q", entry.Redacted)
			}
		}
	}
}

func TestRedact_NilEnvSet(t *testing.T) {
	_, err := Redact(nil, 0)
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestRedact_NegativeSuffix(t *testing.T) {
	es := baseRedactSet()
	_, err := Redact(es, -1)
	if err == nil {
		t.Error("expected error for negative revealSuffix")
	}
}

func TestRedactKeys_SelectedKeys(t *testing.T) {
	es := baseRedactSet()
	out, err := RedactKeys(es, []string{"DATABASE_URL", "API_KEY"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars["DATABASE_URL"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out.Vars["DATABASE_URL"])
	}
	if out.Vars["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out.Vars["API_KEY"])
	}
	if out.Vars["DEBUG"] != "false" {
		t.Errorf("DEBUG should be unchanged, got %q", out.Vars["DEBUG"])
	}
}

func TestRedactKeys_CustomPlaceholder(t *testing.T) {
	es := baseRedactSet()
	out, err := RedactKeys(es, []string{"API_KEY"}, "***HIDDEN***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars["API_KEY"] != "***HIDDEN***" {
		t.Errorf("expected ***HIDDEN***, got %q", out.Vars["API_KEY"])
	}
}

func TestRedactKeys_NilEnvSet(t *testing.T) {
	_, err := RedactKeys(nil, []string{"KEY"}, "")
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestRedactKeys_OriginalUnchanged(t *testing.T) {
	es := baseRedactSet()
	origVal := es.Vars["API_KEY"]
	_, _ = RedactKeys(es, []string{"API_KEY"}, "")
	if es.Vars["API_KEY"] != origVal {
		t.Error("original envset should not be mutated")
	}
}
