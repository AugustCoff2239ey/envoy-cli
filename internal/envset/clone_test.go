package envset_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/envset"
)

func baseCloneSet(t *testing.T) *envset.EnvSet {
	t.Helper()
	es, err := envset.New("myapp", "staging")
	if err != nil {
		t.Fatalf("baseCloneSet: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	return es
}

func TestClone_DefaultNameAndEnv(t *testing.T) {
	src := baseCloneSet(t)
	cloned, err := envset.Clone(src, envset.CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned.Name != src.Name {
		t.Errorf("expected name %q, got %q", src.Name, cloned.Name)
	}
	if cloned.Environment != src.Environment {
		t.Errorf("expected env %q, got %q", src.Environment, cloned.Environment)
	}
	if len(cloned.Values) != len(src.Values) {
		t.Errorf("expected %d values, got %d", len(src.Values), len(cloned.Values))
	}
}

func TestClone_NewNameAndEnv(t *testing.T) {
	src := baseCloneSet(t)
	cloned, err := envset.Clone(src, envset.CloneOptions{
		NewName:        "myapp-copy",
		NewEnvironment: "production",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned.Name != "myapp-copy" {
		t.Errorf("expected name %q, got %q", "myapp-copy", cloned.Name)
	}
	if cloned.Environment != "production" {
		t.Errorf("expected env %q, got %q", "production", cloned.Environment)
	}
}

func TestClone_Independence(t *testing.T) {
	src := baseCloneSet(t)
	cloned, err := envset.Clone(src, envset.CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mutate the clone and verify src is unaffected.
	_ = cloned.Set("DB_HOST", "remotehost")
	if v, _ := src.Get("DB_HOST"); v != "localhost" {
		t.Errorf("src was mutated: expected %q, got %q", "localhost", v)
	}
}

func TestClone_NilSource(t *testing.T) {
	_, err := envset.Clone(nil, envset.CloneOptions{})
	if err == nil {
		t.Error("expected error for nil source, got nil")
	}
}

func TestClone_InvalidNewEnvironment(t *testing.T) {
	src := baseCloneSet(t)
	_, err := envset.Clone(src, envset.CloneOptions{NewEnvironment: "bad env!"})
	if err == nil {
		t.Error("expected error for invalid environment, got nil")
	}
}
