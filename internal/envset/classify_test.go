package envset

import (
	"testing"
)

func baseClassifySet() *EnvSet {
	es, _ := New("classify-test", "local")
	_ = es.Set("DB_PASSWORD", "secret123")
	_ = es.Set("DATABASE_URL", "postgres://localhost/mydb")
	_ = es.Set("API_TOKEN", "tok_abc")
	_ = es.Set("SERVER_HOST", "localhost")
	_ = es.Set("FEATURE_FLAG_DARK_MODE", "true")
	_ = es.Set("APP_NAME", "myapp")
	return es
}

func TestClassify_Categories(t *testing.T) {
	es := baseClassifySet()
	report, err := Classify(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 6 {
		t.Errorf("expected 6 results, got %d", len(report.Results))
	}
	byKey := make(map[string]string)
	for _, r := range report.Results {
		byKey[r.Key] = r.Category
	}
	if byKey["DB_PASSWORD"] != "secret" {
		t.Errorf("DB_PASSWORD: want secret, got %s", byKey["DB_PASSWORD"])
	}
	if byKey["DATABASE_URL"] != "database" {
		t.Errorf("DATABASE_URL: want database, got %s", byKey["DATABASE_URL"])
	}
	if byKey["API_TOKEN"] != "secret" {
		t.Errorf("API_TOKEN: want secret, got %s", byKey["API_TOKEN"])
	}
	if byKey["SERVER_HOST"] != "network" {
		t.Errorf("SERVER_HOST: want network, got %s", byKey["SERVER_HOST"])
	}
	if byKey["FEATURE_FLAG_DARK_MODE"] != "feature" {
		t.Errorf("FEATURE_FLAG_DARK_MODE: want feature, got %s", byKey["FEATURE_FLAG_DARK_MODE"])
	}
	if byKey["APP_NAME"] != "general" {
		t.Errorf("APP_NAME: want general, got %s", byKey["APP_NAME"])
	}
}

func TestClassify_ByCategory(t *testing.T) {
	es := baseClassifySet()
	report, _ := Classify(es)
	if len(report.ByCategory["secret"]) < 1 {
		t.Error("expected at least one secret key")
	}
	if len(report.ByCategory["general"]) < 1 {
		t.Error("expected at least one general key")
	}
}

func TestClassify_NilEnvSet(t *testing.T) {
	_, err := Classify(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestClassify_EmptySet(t *testing.T) {
	es, _ := New("empty", "local")
	report, err := Classify(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(report.Results))
	}
}
