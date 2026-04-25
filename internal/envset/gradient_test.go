package envset

import (
	"testing"
)

func baseGradientSets(t *testing.T) []*EnvSet {
	t.Helper()
	local, _ := New("app", "local")
	_ = local.Set("DB_HOST", "localhost")
	_ = local.Set("LOG_LEVEL", "debug")
	_ = local.Set("SHARED", "same")

	staging, _ := New("app", "staging")
	_ = staging.Set("DB_HOST", "staging.db.internal")
	_ = staging.Set("LOG_LEVEL", "info")
	_ = staging.Set("SHARED", "same")

	prod, _ := New("app", "production")
	_ = prod.Set("DB_HOST", "prod.db.internal")
	_ = prod.Set("LOG_LEVEL", "warn")
	_ = prod.Set("SHARED", "same")

	return []*EnvSet{local, staging, prod}
}

func TestGradient_AllKeys(t *testing.T) {
	sets := baseGradientSets(t)
	results, err := Gradient(sets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
}

func TestGradient_SelectedKey(t *testing.T) {
	sets := baseGradientSets(t)
	results, err := Gradient(sets, []string{"DB_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", results[0].Key)
	}
	if results[0].Uniform {
		t.Error("expected DB_HOST to not be uniform across envs")
	}
	if len(results[0].Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(results[0].Steps))
	}
}

func TestGradient_UniformKey(t *testing.T) {
	sets := baseGradientSets(t)
	results, err := Gradient(sets, []string{"SHARED"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Uniform {
		t.Error("expected SHARED to be uniform across envs")
	}
}

func TestGradient_NoSets(t *testing.T) {
	_, err := Gradient(nil, nil)
	if err == nil {
		t.Error("expected error for nil sets")
	}
}

func TestGradient_InvalidKey(t *testing.T) {
	sets := baseGradientSets(t)
	_, err := Gradient(sets, []string{"invalid key!"})
	if err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestGradient_MissingKeyInSomeEnvs(t *testing.T) {
	sets := baseGradientSets(t)
	// ONLY_LOCAL is only in local
	_ = sets[0].Set("ONLY_LOCAL", "value")
	results, err := Gradient(sets, []string{"ONLY_LOCAL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Uniform {
		t.Error("expected non-uniform result when key is missing in some envs")
	}
	if results[0].Steps[1].Value != "" {
		t.Errorf("expected empty value for missing key in staging, got %q", results[0].Steps[1].Value)
	}
}
