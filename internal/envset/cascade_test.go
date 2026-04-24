package envset_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/envset"
)

func baseCascadeSets(t *testing.T) (*envset.EnvSet, *envset.EnvSet, *envset.EnvSet) {
	t.Helper()

	base, err := envset.New("base", "production")
	if err != nil {
		t.Fatalf("failed to create base set: %v", err)
	}
	_ = base.Set("DB_HOST", "prod-db.internal")
	_ = base.Set("DB_PORT", "5432")
	_ = base.Set("LOG_LEVEL", "warn")
	_ = base.Set("FEATURE_X", "false")

	mid, err := envset.New("mid", "staging")
	if err != nil {
		t.Fatalf("failed to create mid set: %v", err)
	}
	_ = mid.Set("DB_HOST", "staging-db.internal")
	_ = mid.Set("LOG_LEVEL", "info")

	override, err := envset.New("override", "local")
	if err != nil {
		t.Fatalf("failed to create override set: %v", err)
	}
	_ = override.Set("LOG_LEVEL", "debug")
	_ = override.Set("FEATURE_X", "true")

	return base, mid, override
}

func TestCascade_BasicPrecedence(t *testing.T) {
	base, mid, override := baseCascadeSets(t)

	result, err := envset.Cascade(envset.DefaultCascadeOptions(), base, mid, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// override wins for LOG_LEVEL
	if v, _ := result.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", v)
	}

	// mid wins for DB_HOST (override doesn't set it)
	if v, _ := result.Get("DB_HOST"); v != "staging-db.internal" {
		t.Errorf("expected DB_HOST=staging-db.internal, got %q", v)
	}

	// base supplies DB_PORT (no one else sets it)
	if v, _ := result.Get("DB_PORT"); v != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", v)
	}

	// override wins for FEATURE_X
	if v, _ := result.Get("FEATURE_X"); v != "true" {
		t.Errorf("expected FEATURE_X=true, got %q", v)
	}
}

func TestCascade_SingleLayer(t *testing.T) {
	base, _, _ := baseCascadeSets(t)

	result, err := envset.Cascade(envset.DefaultCascadeOptions(), base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v, _ := result.Get("DB_HOST"); v != "prod-db.internal" {
		t.Errorf("expected DB_HOST=prod-db.internal, got %q", v)
	}
}

func TestCascade_NoLayers(t *testing.T) {
	_, err := envset.Cascade(envset.DefaultCascadeOptions())
	if err == nil {
		t.Error("expected error for empty layer list, got nil")
	}
}

func TestCascade_NilLayerSkipped(t *testing.T) {
	base, _, override := baseCascadeSets(t)

	result, err := envset.Cascade(envset.DefaultCascadeOptions(), base, nil, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// override still wins for LOG_LEVEL
	if v, _ := result.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", v)
	}

	// base supplies DB_HOST (mid was nil)
	if v, _ := result.Get("DB_HOST"); v != "prod-db.internal" {
		t.Errorf("expected DB_HOST=prod-db.internal, got %q", v)
	}
}

func TestCascade_ResultIndependence(t *testing.T) {
	base, mid, override := baseCascadeSets(t)

	result, err := envset.Cascade(envset.DefaultCascadeOptions(), base, mid, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// mutating source should not affect result
	_ = override.Set("LOG_LEVEL", "error")

	if v, _ := result.Get("LOG_LEVEL"); v != "debug" {
		t.Errorf("cascade result was mutated by source change: got %q", v)
	}
}
