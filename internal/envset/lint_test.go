package envset

import (
	"strings"
	"testing"
)

func baseLintSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("lint-test", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["API_KEY"] = "abc123"
	es.Vars["DB_HOST"] = "localhost"
	return es
}

func TestLint_CleanSet(t *testing.T) {
	es := baseLintSet(t)
	res, err := Lint(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Findings) != 0 {
		t.Errorf("expected no findings, got: %s", res.Summary())
	}
}

func TestLint_EmptyValue(t *testing.T) {
	es := baseLintSet(t)
	es.Vars["EMPTY_VAR"] = ""
	res, err := Lint(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range res.Findings {
		if f.Key == "EMPTY_VAR" && f.Severity == LintWarn {
			found = true
		}
	}
	if !found {
		t.Error("expected warn finding for empty value")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	es := baseLintSet(t)
	es.Vars["lowercase_key"] = "value"
	res, err := Lint(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasErrors() {
		t.Error("expected at least one error finding for lowercase key")
	}
}

func TestLint_ShellBuiltin(t *testing.T) {
	es := baseLintSet(t)
	es.Vars["PATH"] = "/usr/local/bin"
	res, err := Lint(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range res.Findings {
		if f.Key == "PATH" && strings.Contains(f.Message, "shell built-in") {
			found = true
		}
	}
	if !found {
		t.Error("expected warn finding for shell built-in key")
	}
}

func TestLint_NilEnvSet(t *testing.T) {
	_, err := Lint(nil)
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestLint_Summary_NoIssues(t *testing.T) {
	es := baseLintSet(t)
	res, _ := Lint(es)
	if res.Summary() != "no issues found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
