package envset

import (
	"testing"
)

func baseBaselineSet() *EnvSet {
	e, _ := New("baseline-test", "staging")
	_ = e.Set("DB_HOST", "localhost")
	_ = e.Set("DB_PORT", "5432")
	_ = e.Set("APP_ENV", "staging")
	return e
}

func TestSetBaseline_Valid(t *testing.T) {
	e := baseBaselineSet()
	b, err := SetBaseline(e, "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Label != "v1" {
		t.Errorf("expected label v1, got %s", b.Label)
	}
	if len(b.Vars) != 3 {
		t.Errorf("expected 3 vars, got %d", len(b.Vars))
	}
}

func TestSetBaseline_NilEnvSet(t *testing.T) {
	_, err := SetBaseline(nil, "v1")
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestSetBaseline_EmptyLabel(t *testing.T) {
	e := baseBaselineSet()
	_, err := SetBaseline(e, "")
	if err == nil {
		t.Error("expected error for empty label")
	}
}

func TestSetBaseline_Independence(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "snap")
	_ = e.Set("DB_HOST", "remotehost")
	if b.Vars["DB_HOST"] != "localhost" {
		t.Error("baseline should not reflect mutations to source envset")
	}
}

func TestDetectDrift_NoDrift(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "clean")
	result, err := DetectDrift(e, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestDetectDrift_Added(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "pre")
	_ = e.Set("NEW_KEY", "value")
	result, _ := DetectDrift(e, b)
	if _, ok := result.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY in added")
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "pre")
	delete(e.Vars, "APP_ENV")
	result, _ := DetectDrift(e, b)
	if _, ok := result.Removed["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in removed")
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "pre")
	_ = e.Set("DB_PORT", "3306")
	result, _ := DetectDrift(e, b)
	pair, ok := result.Changed["DB_PORT"]
	if !ok {
		t.Fatal("expected DB_PORT in changed")
	}
	if pair[0] != "5432" || pair[1] != "3306" {
		t.Errorf("unexpected change pair: %v", pair)
	}
}

func TestDetectDrift_NilInputs(t *testing.T) {
	e := baseBaselineSet()
	b, _ := SetBaseline(e, "b")
	if _, err := DetectDrift(nil, b); err == nil {
		t.Error("expected error for nil envset")
	}
	if _, err := DetectDrift(e, nil); err == nil {
		t.Error("expected error for nil baseline")
	}
}
