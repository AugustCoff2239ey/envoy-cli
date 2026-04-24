package envset

import (
	"strings"
	"testing"
)

func baseTruncateSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("trunc", "test")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["SHORT"] = "hi"
	es.Vars["EXACT"] = strings.Repeat("x", 64)
	es.Vars["LONG"] = strings.Repeat("a", 100)
	es.Vars["VERY_LONG"] = strings.Repeat("b", 200)
	return es
}

func TestTruncate_AllKeys(t *testing.T) {
	es := baseTruncateSet(t)
	opts := DefaultTruncateOptions()

	out, err := Truncate(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged, got %q", out["SHORT"])
	}
	if out["EXACT"] != strings.Repeat("x", 64) {
		t.Errorf("EXACT should be unchanged")
	}
	if len(out["LONG"]) != 64 {
		t.Errorf("LONG should be truncated to 64 chars, got %d", len(out["LONG"]))
	}
	if !strings.HasSuffix(out["LONG"], "...") {
		t.Errorf("LONG should end with '...', got %q", out["LONG"])
	}
}

func TestTruncate_SelectedKeys(t *testing.T) {
	es := baseTruncateSet(t)
	opts := DefaultTruncateOptions()
	opts.Keys = []string{"LONG"}

	out, err := Truncate(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["VERY_LONG"]; ok {
		t.Errorf("VERY_LONG should not appear in output when not in Keys")
	}
	if len(out["LONG"]) != 64 {
		t.Errorf("LONG should be truncated, got len %d", len(out["LONG"]))
	}
}

func TestTruncate_CustomSuffix(t *testing.T) {
	es := baseTruncateSet(t)
	opts := TruncateOptions{MaxLength: 10, Suffix: "~"}

	out, err := Truncate(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(out["LONG"], "~") {
		t.Errorf("LONG should end with '~', got %q", out["LONG"])
	}
	if len(out["LONG"]) != 10 {
		t.Errorf("expected length 10, got %d", len(out["LONG"]))
	}
}

func TestTruncate_NilEnvSet(t *testing.T) {
	_, err := Truncate(nil, DefaultTruncateOptions())
	if err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestTruncate_InvalidMaxLength(t *testing.T) {
	es := baseTruncateSet(t)
	_, err := Truncate(es, TruncateOptions{MaxLength: 0, Suffix: "..."})
	if err == nil {
		t.Fatal("expected error for MaxLength=0")
	}
}

func TestTruncateApply_ModifiesInPlace(t *testing.T) {
	es := baseTruncateSet(t)
	opts := DefaultTruncateOptions()

	if err := TruncateApply(es, opts); err != nil {
		t.Fatalf("TruncateApply: %v", err)
	}
	if len(es.Vars["LONG"]) != 64 {
		t.Errorf("expected LONG to be 64 chars in-place, got %d", len(es.Vars["LONG"]))
	}
	if es.Vars["SHORT"] != "hi" {
		t.Errorf("SHORT should remain unchanged")
	}
}
