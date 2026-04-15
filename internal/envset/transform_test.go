package envset

import (
	"strings"
	"testing"
)

func baseTransformSet() *EnvSet {
	es, _ := New("transform-test	es := baseTransformSet()
	if err := Transform(es, "upper", nil, TransformOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["APP_ENV"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", es.Vars["APP_ENV"])
	}
}

func TestTransform_Lower(t *testing.T) {
	es := baseTransformSet()
	es.Vars["APP_ENV"] = "STAGING"
	if err := Transform(es, "lower", nil, TransformOptions{Keys: []string{"APP_ENV"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["APP_ENV"] != "staging" {
		t.Errorf("expected staging, got %q", es.Vars["APP_ENV"])
	}
}

func TestTransform_Trim(t *testing.T) {
	es := baseTransformSet()
	if err := Transform(es, "trim", nil, TransformOptions{Keys: []string{"APP_NAME"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["APP_NAME"] != "myapp" {
		t.Errorf("expected myapp, got %q", es.Vars["APP_NAME"])
	}
}

func TestTransform_CustomFunc(t *testing.T) {
	es := baseTransformSet()
	doubler := func(v string) (string, error) { return v + v, nil }
	if err := Transform(es, "", doubler, TransformOptions{Keys: []string{"DB_HOST"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["DB_HOST"] != "localhostlocalhost" {
		t.Errorf("expected localhostlocalhost, got %q", es.Vars["DB_HOST"])
	}
}

func TestTransform_UnknownBuiltin(t *testing.T) {
	es := baseTransformSet()
	err := Transform(es, "nonexistent", nil, TransformOptions{})
	if err == nil || !strings.Contains(err.Error(), "unknown transform") {
		t.Errorf("expected unknown transform error, got %v", err)
	}
}

func TestTransform_MissingKey(t *testing.T) {
	es := baseTransformSet()
	err := Transform(es, "upper", nil, TransformOptions{Keys: []string{"MISSING_KEY"}})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestTransform_SkipLocked(t *testing.T) {
	es := baseTransformSet()
	_ = LockKey(es, "APP_ENV")
	origVal := es.Vars["APP_ENV"]
	if err := Transform(es, "upper", nil, TransformOptions{SkipLocked: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["APP_ENV"] != origVal {
		t.Errorf("locked key should not be transformed")
	}
}

func TestTransform_NilEnvSet(t *testing.T) {
	err := Transform(nil, "upper", nil, TransformOptions{})
	if err == nil {
		t.Error("expected error for nil envset")
	}
}
