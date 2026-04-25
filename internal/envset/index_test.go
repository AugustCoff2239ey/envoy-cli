package envset

import (
	"testing"
)

func baseIndexSet() *EnvSet {
	es, _ := New("indexset", "test")
	_ = es.Set("APP_HOST", "localhost")
	_ = es.Set("APP_PORT", "8080")
	_ = es.Set("DB_URL", "postgres://localhost/db")
	es.Meta["group:APP_HOST"] = "app"
	es.Meta["group:APP_PORT"] = "app"
	es.Meta["group:DB_URL"] = "database"
	es.Meta["tags:APP_HOST"] = "core,required"
	return es
}

func TestIndex_KeysPresent(t *testing.T) {
	es := baseIndexSet()
	idx, err := Index(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"APP_HOST", "APP_PORT", "DB_URL"} {
		if _, ok := idx[k]; !ok {
			t.Errorf("expected key %q in index", k)
		}
	}
}

func TestIndex_PositionsAreUnique(t *testing.T) {
	es := baseIndexSet()
	idx, _ := Index(es)
	seen := map[int]string{}
	for k, e := range idx {
		if prev, dup := seen[e.Position]; dup {
			t.Errorf("position %d shared by %q and %q", e.Position, prev, k)
		}
		seen[e.Position] = k
	}
}

func TestIndex_GroupAssigned(t *testing.T) {
	es := baseIndexSet()
	idx, _ := Index(es)
	if idx["APP_HOST"].Group != "app" {
		t.Errorf("expected group 'app' for APP_HOST, got %q", idx["APP_HOST"].Group)
	}
	if idx["DB_URL"].Group != "database" {
		t.Errorf("expected group 'database' for DB_URL, got %q", idx["DB_URL"].Group)
	}
}

func TestIndex_TagsParsed(t *testing.T) {
	es := baseIndexSet()
	idx, _ := Index(es)
	tags := idx["APP_HOST"].Tags
	if len(tags) != 2 || tags[0] != "core" || tags[1] != "required" {
		t.Errorf("unexpected tags for APP_HOST: %v", tags)
	}
}

func TestIndex_NilEnvSet(t *testing.T) {
	_, err := Index(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestIndexByGroup_Filtered(t *testing.T) {
	es := baseIndexSet()
	idx, _ := Index(es)
	entries := IndexByGroup(idx, "app")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in 'app' group, got %d", len(entries))
	}
}

func TestIndexByGroup_Empty(t *testing.T) {
	es := baseIndexSet()
	idx, _ := Index(es)
	entries := IndexByGroup(idx, "nonexistent")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for unknown group, got %d", len(entries))
	}
}
