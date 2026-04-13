package envset

import (
	"testing"
)

func baseResolveSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("resolve-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("APP_HOST", "localhost")
	_ = es.Set("APP_PORT", "8080")
	return es
}

func TestResolve_LocalSource(t *testing.T) {
	es := baseResolveSet(t)
	results, err := Resolve(es, ResolveOptions{})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	for _, r := range results {
		if r.Source != "local" {
			t.Errorf("expected source=local for key %q, got %q", r.Key, r.Source)
		}
	}
}

func TestResolve_OverrideTakesPriority(t *testing.T) {
	es := baseResolveSet(t)
	opts := ResolveOptions{
		Overrides: map[string]string{"APP_PORT": "9090"},
	}
	results, err := Resolve(es, opts)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	for _, r := range results {
		if r.Key == "APP_PORT" {
			if r.Resolved != "9090" || r.Source != "override" {
				t.Errorf("expected override value 9090/override, got %q/%q", r.Resolved, r.Source)
			}
		}
	}
}

func TestResolve_DefaultFallback(t *testing.T) {
	es := baseResolveSet(t)
	opts := ResolveOptions{
		Defaults: map[string]string{"APP_DEBUG": "false"},
	}
	results, err := Resolve(es, opts)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	found := false
	for _, r := range results {
		if r.Key == "APP_DEBUG" {
			found = true
			if r.Resolved != "false" || r.Source != "default" {
				t.Errorf("expected default value false/default, got %q/%q", r.Resolved, r.Source)
			}
		}
	}
	if !found {
		t.Error("expected APP_DEBUG in results via default")
	}
}

func TestResolve_NilEnvSet(t *testing.T) {
	_, err := Resolve(nil, ResolveOptions{})
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestResolveKey_Override(t *testing.T) {
	es := baseResolveSet(t)
	opts := ResolveOptions{Overrides: map[string]string{"APP_HOST": "prod.example.com"}}
	r, err := ResolveKey(es, "APP_HOST", opts)
	if err != nil {
		t.Fatalf("ResolveKey: %v", err)
	}
	if r.Resolved != "prod.example.com" || r.Source != "override" {
		t.Errorf("unexpected result: %+v", r)
	}
}

func TestResolveKey_NotFound(t *testing.T) {
	es := baseResolveSet(t)
	_, err := ResolveKey(es, "MISSING_KEY", ResolveOptions{})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestResolveKey_NilEnvSet(t *testing.T) {
	_, err := ResolveKey(nil, "APP_HOST", ResolveOptions{})
	if err == nil {
		t.Error("expected error for nil envset")
	}
}
