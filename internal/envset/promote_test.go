package envset_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/envset"
)

func basePromoteSets(t *testing.T) (source, target *envset.EnvSet) {
	t.Helper()
	source, _ = envset.New("app", "staging")
	_ = source.Set("DB_HOST", "staging-db.internal")
	_ = source.Set("API_KEY", "stg-secret")
	_ = source.Set("LOG_LEVEL", "debug")

	target, _ = envset.New("app", "production")
	_ = target.Set("DB_HOST", "prod-db.internal")
	_ = target.Set("LOG_LEVEL", "warn")
	return
}

func TestPromote_AllKeys_NoOverwrite(t *testing.T) {
	source, target := basePromoteSets(t)
	result, err := envset.Promote(source, target, envset.PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Existing keys must NOT be overwritten.
	if result.Vars["DB_HOST"] != "prod-db.internal" {
		t.Errorf("DB_HOST should not be overwritten, got %q", result.Vars["DB_HOST"])
	}
	// New key from source should be present.
	if result.Vars["API_KEY"] != "stg-secret" {
		t.Errorf("API_KEY should be promoted, got %q", result.Vars["API_KEY"])
	}
}

func TestPromote_AllKeys_WithOverwrite(t *testing.T) {
	source, target := basePromoteSets(t)
	result, err := envset.Promote(source, target, envset.PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Vars["DB_HOST"] != "staging-db.internal" {
		t.Errorf("DB_HOST should be overwritten, got %q", result.Vars["DB_HOST"])
	}
	if result.Vars["LOG_LEVEL"] != "debug" {
		t.Errorf("LOG_LEVEL should be overwritten, got %q", result.Vars["LOG_LEVEL"])
	}
}

func TestPromote_SelectedKeys(t *testing.T) {
	source, target := basePromoteSets(t)
	result, err := envset.Promote(source, target, envset.PromoteOptions{
		Keys:      []string{"API_KEY"},
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Vars["API_KEY"] != "stg-secret" {
		t.Errorf("expected API_KEY to be promoted")
	}
	// DB_HOST should remain the target value even though Overwrite is true,
	// because it was not in the selected keys list.
	if result.Vars["DB_HOST"] != "prod-db.internal" {
		t.Errorf("DB_HOST should not be touched, got %q", result.Vars["DB_HOST"])
	}
}

func TestPromote_MissingKey(t *testing.T) {
	source, target := basePromoteSets(t)
	_, err := envset.Promote(source, target, envset.PromoteOptions{Keys: []string{"NONEXISTENT"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPromote_NilSource(t *testing.T) {
	_, target := basePromoteSets(t)
	_, err := envset.Promote(nil, target, envset.PromoteOptions{})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestPromote_TargetEnvironmentPreserved(t *testing.T) {
	source, target := basePromoteSets(t)
	result, err := envset.Promote(source, target, envset.PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", result.Environment)
	}
}
