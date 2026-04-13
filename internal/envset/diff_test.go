package envset

import (
	"strings"
	"testing"
)

func baseAndTarget(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	base, _ := New("base", "local")
	base.Vars["FOO"] = "foo"
	base.Vars["BAR"] = "bar"
	base.Vars["SHARED"] = "same"

	target, _ := New("target", "staging")
	target.Vars["BAZ"] = "baz"
	target.Vars["SHARED"] = "same"
	target.Vars["FOO"] = "changed"

	return base, target
}

func TestDiff_Added(t *testing.T) {
	base, target := baseAndTarget(t)
	result := Diff(base, target)
	if _, ok := result.Added["BAZ"]; !ok {
		t.Error("expected BAZ to be in Added")
	}
}

func TestDiff_Removed(t *testing.T) {
	base, target := baseAndTarget(t)
	result := Diff(base, target)
	if _, ok := result.Removed["BAR"]; !ok {
		t.Error("expected BAR to be in Removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	base, target := baseAndTarget(t)
	result := Diff(base, target)
	if vals, ok := result.Changed["FOO"]; !ok {
		t.Error("expected FOO to be in Changed")
	} else if vals[0] != "foo" || vals[1] != "changed" {
		t.Errorf("unexpected changed values: %v", vals)
	}
}

func TestDiff_NoDiff(t *testing.T) {
	a, _ := New("a", "local")
	a.Vars["KEY"] = "val"
	b, _ := New("b", "local")
	b.Vars["KEY"] = "val"

	result := Diff(a, b)
	if result.HasDiff() {
		t.Error("expected no diff")
	}
	if result.String() != "No differences found." {
		t.Errorf("unexpected string: %s", result.String())
	}
}

func TestDiff_StringContainsMarkers(t *testing.T) {
	base, target := baseAndTarget(t)
	result := Diff(base, target)
	s := result.String()
	if !strings.Contains(s, "+ BAZ") {
		t.Error("expected '+ BAZ' in output")
	}
	if !strings.Contains(s, "- BAR") {
		t.Error("expected '- BAR' in output")
	}
	if !strings.Contains(s, "~ FOO") {
		t.Error("expected '~ FOO' in output")
	}
}
