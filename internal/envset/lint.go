package envset

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint finding.
type LintSeverity string

const (
	LintWarn  LintSeverity = "warn"
	LintError LintSeverity = "error"
)

// LintFinding represents a single lint result for a key.
type LintFinding struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintResult holds all findings from a lint run.
type LintResult struct {
	Findings []LintFinding
}

// HasErrors returns true if any finding has error severity.
func (r *LintResult) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == LintError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *LintResult) Summary() string {
	if len(r.Findings) == 0 {
		return "no issues found"
	}
	var sb strings.Builder
	for _, f := range r.Findings {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", f.Severity, f.Key, f.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Lint inspects an EnvSet for common issues and returns a LintResult.
// It checks for empty values, overly long values, lowercase keys, and
// keys that shadow common shell variables.
func Lint(es *EnvSet) (*LintResult, error) {
	if es == nil {
		return nil, fmt.Errorf("lint: envset must not be nil")
	}

	result := &LintResult{}

	shellBuiltins := map[string]bool{
		"PATH": true, "HOME": true, "USER": true, "SHELL": true,
		"TERM": true, "LANG": true, "PWD": true, "OLDPWD": true,
	}

	for key, val := range es.Vars {
		if val == "" {
			result.Findings = append(result.Findings, LintFinding{
				Key:      key,
				Message:  "value is empty",
				Severity: LintWarn,
			})
		}
		if len(val) > 1024 {
			result.Findings = append(result.Findings, LintFinding{
				Key:      key,
				Message:  fmt.Sprintf("value exceeds 1024 characters (%d)", len(val)),
				Severity: LintWarn,
			})
		}
		if key != strings.ToUpper(key) {
			result.Findings = append(result.Findings, LintFinding{
				Key:      key,
				Message:  "key should be uppercase",
				Severity: LintError,
			})
		}
		if shellBuiltins[strings.ToUpper(key)] {
			result.Findings = append(result.Findings, LintFinding{
				Key:      key,
				Message:  "key shadows a common shell built-in variable",
				Severity: LintWarn,
			})
		}
	}

	return result, nil
}
