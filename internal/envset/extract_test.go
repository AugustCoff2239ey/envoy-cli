package envset

import (
	"testing"
)

func baseExtractSet(t *testing.T) *EnvSet {
	t.Helper()
	es, _ := New("myapp", "staging")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("APP_NAME", "envoy")
	_ = es.Set("APP_ENV", "staging")
	_ = es.Set("SECRET_KEY", "abc123")
	return es
}

func TestExtract_ByPattern(t *testing.T) {
	es := baseExtractSet(t)
	res, err := Extract(es, ExtractOptions{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Keys))
	}
	if _, ok := res.Extracted.Vars["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in extracted set")
	}
}

func TestExtract_ByKeys(t *testing.T) {
	es := baseExtractSet(t)
	res, err := Extract(es, ExtractOptions{Keys: []string{"APP_NAME", "SECRET_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Keys))
	}
	if res.Extracted.Vars["APP_NAME"] != "envoy" {
		t.Errorf("unexpected value for APP_NAME")
	}
}

func TestExtract_RemoveFromSource(t *testing.T) {
	es := baseExtractSet(t)
	_, err := Extract(es, ExtractOptions{Pattern: "^DB_", RemoveFromSource: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := es.Vars["DB_HOST"]; ok {
		t.Error("DB_HOST should have been removed from source")
	}
	if _, ok := es.Vars["APP_NAME"]; !ok {
		t.Error("APP_NAME should remain in source")
	}
}

func TestExtract_MissingKey(t *testing.T) {
	es := baseExtractSet(t)
	_, err := Extract(es, ExtractOptions{Keys: []string{"MISSING_KEY"}})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestExtract_InvalidPattern(t *testing.T) {
	es := baseExtractSet(t)
	_, err := Extract(es, ExtractOptions{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestExtract_NilEnvSet(t *testing.T) {
	_, err := Extract(nil, ExtractOptions{Pattern: "^DB_"})
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestExtract_EmptyOptions(t *testing.T) {
	es := baseExtractSet(t)
	res, err := Extract(es, ExtractOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 0 {
		t.Errorf("expected 0 keys with empty options, got %d", len(res.Keys))
	}
}
