package envset

import (
	"testing"
)

func baseCompareSets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	base, _ := New("app", "staging")
	base.Vars["HOST"] = "localhost"
	base.Vars["PORT"] = "8080"
	base.Vars["ONLY_BASE"] = "yes"

	target, _ := New("app", "production")
	target.Vars["HOST"] = "example.com"
	target.Vars["PORT"] = "8080"
	target.Vars["ONLY_TARGET"] = "yes"

	return base, target
}

func TestCompare_Matching(t *testing.T) {
	base, target := baseCompareSets(t)
	res, err := Compare(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.MatchingKeys) != 1 || res.MatchingKeys[0] != "PORT" {
		t.Errorf("expected MatchingKeys=[PORT], got %v", res.MatchingKeys)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	base, target := baseCompareSets(t)
	res, _ := Compare(base, target)
	pair, ok := res.MismatchedKeys["HOST"]
	if !ok {
		t.Fatal("expected HOST in MismatchedKeys")
	}
	if pair[0] != "localhost" || pair[1] != "example.com" {
		t.Errorf("unexpected mismatch values: %v", pair)
	}
}

func TestCompare_OnlyInBase(t *testing.T) {
	base, target := baseCompareSets(t)
	res, _ := Compare(base, target)
	if len(res.OnlyInBase) != 1 || res.OnlyInBase[0] != "ONLY_BASE" {
		t.Errorf("expected OnlyInBase=[ONLY_BASE], got %v", res.OnlyInBase)
	}
}

func TestCompare_OnlyInTarget(t *testing.T) {
	base, target := baseCompareSets(t)
	res, _ := Compare(base, target)
	if len(res.OnlyInTarget) != 1 || res.OnlyInTarget[0] != "ONLY_TARGET" {
		t.Errorf("expected OnlyInTarget=[ONLY_TARGET], got %v", res.OnlyInTarget)
	}
}

func TestCompare_Equal(t *testing.T) {
	base, _ := New("app", "staging")
	base.Vars["KEY"] = "value"
	target, _ := New("app", "production")
	target.Vars["KEY"] = "value"
	res, _ := Compare(base, target)
	if !res.Equal {
		t.Error("expected Equal=true for identical var sets")
	}
}

func TestCompare_NilInputs(t *testing.T) {
	base, _ := New("app", "staging")
	_, err := Compare(base, nil)
	if err == nil {
		t.Error("expected error for nil target")
	}
	_, err = Compare(nil, base)
	if err == nil {
		t.Error("expected error for nil base")
	}
}
