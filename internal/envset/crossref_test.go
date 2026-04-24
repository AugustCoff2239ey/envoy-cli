package envset

import (
	"testing"
)

func baseCrossRefSets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	base, _ := New("base", "staging")
	base.Vars["DB_HOST"] = "localhost"
	base.Vars["DB_PORT"] = "5432"
	base.Vars["API_KEY"] = "secret"
	base.Vars["BASE_ONLY"] = "only-in-base"

	target, _ := New("target", "production")
	target.Vars["DB_HOST"] = "prod-host"
	target.Vars["DB_PORT"] = "5432"
	target.Vars["API_KEY"] = "prod-secret"
	target.Vars["TARGET_ONLY"] = "only-in-target"

	return base, target
}

func TestCrossRef_SharedKeys(t *testing.T) {
	base, target := baseCrossRefSets(t)
	res, err := CrossRef(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.SharedKeys) != 3 {
		t.Errorf("expected 3 shared keys, got %d", len(res.SharedKeys))
	}
}

func TestCrossRef_OnlyInBase(t *testing.T) {
	base, target := baseCrossRefSets(t)
	res, _ := CrossRef(base, target)
	if len(res.OnlyInBase) != 1 || res.OnlyInBase[0] != "BASE_ONLY" {
		t.Errorf("expected [BASE_ONLY], got %v", res.OnlyInBase)
	}
}

func TestCrossRef_OnlyInTarget(t *testing.T) {
	base, target := baseCrossRefSets(t)
	res, _ := CrossRef(base, target)
	if len(res.OnlyInTarget) != 1 || res.OnlyInTarget[0] != "TARGET_ONLY" {
		t.Errorf("expected [TARGET_ONLY], got %v", res.OnlyInTarget)
	}
}

func TestCrossRef_ValueMatches(t *testing.T) {
	base, target := baseCrossRefSets(t)
	res, _ := CrossRef(base, target)
	if len(res.ValueMatches) != 1 || res.ValueMatches[0] != "DB_PORT" {
		t.Errorf("expected [DB_PORT], got %v", res.ValueMatches)
	}
}

func TestCrossRef_ValueMismatches(t *testing.T) {
	base, target := baseCrossRefSets(t)
	res, _ := CrossRef(base, target)
	if len(res.ValueMismatches) != 2 {
		t.Errorf("expected 2 mismatches, got %d: %v", len(res.ValueMismatches), res.ValueMismatches)
	}
}

func TestCrossRef_NilInputs(t *testing.T) {
	base, _ := baseCrossRefSets(t)
	_, err := CrossRef(base, nil)
	if err == nil {
		t.Error("expected error for nil target")
	}
	_, err = CrossRef(nil, base)
	if err == nil {
		t.Error("expected error for nil base")
	}
}
