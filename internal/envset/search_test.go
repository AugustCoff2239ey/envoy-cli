package envset

import (
	"testing"
)

func baseSearchSet(t *testing.T) *EnvSet {
	t.Helper()
	es, _ := New("search-test", "local")
	es.Vars["DATABASE_URL"] = "postgres://localhost/mydb"
	es.Vars["REDIS_URL"] = "redis://localhost:6379"
	es.Vars["APP_SECRET"] = "supersecret"
	es.Vars["LOG_LEVEL"] = "debug"
	es.Vars["DB_POOL_SIZE"] = "10"
	return es
}

func TestSearch_ByKeySubstring(t *testing.T) {
	es := baseSearchSet(t)
	results, err := Search(es, "DB", SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ByValueSubstring(t *testing.T) {
	es := baseSearchSet(t)
	results, err := Search(es, "localhost", SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_KeysOnly(t *testing.T) {
	es := baseSearchSet(t)
	results, err := Search(es, "url", SearchOptions{KeysOnly: true, CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Field != "key" {
			t.Errorf("expected field=key, got %q", r.Field)
		}
	}
}

func TestSearch_Regex(t *testing.T) {
	es := baseSearchSet(t)
	results, err := Search(es, "^(DATABASE|REDIS)_URL$", SearchOptions{Regex: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestSearch_InvalidRegex(t *testing.T) {
	es := baseSearchSet(t)
	_, err := Search(es, "[invalid", SearchOptions{Regex: true})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	es := baseSearchSet(t)
	_, err := Search(es, "", SearchOptions{})
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestSearch_NilEnvSet(t *testing.T) {
	_, err := Search(nil, "DB", SearchOptions{})
	if err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestSearch_NoResults(t *testing.T) {
	es := baseSearchSet(t)
	results, err := Search(es, "NONEXISTENT_XYZ", SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_CaseSensitive(t *testing.T) {
	es := baseSearchSet(t)
	// With case-sensitive search, lowercase "db" should not match "DB_POOL_SIZE" or "DATABASE_URL"
	results, err := Search(es, "db", SearchOptions{CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for case-sensitive search, got %d", len(results))
	}
}
