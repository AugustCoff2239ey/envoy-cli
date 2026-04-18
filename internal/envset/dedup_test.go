package envset

import (
	"testing"
)

func baseDedupSets(t *testing.T) (*EnvSet, *EnvSet, *EnvSet) {
	t.Helper()
	base, _ := New("app", "local")
	base.Vars["DB_HOST"] = "localhost"
	base.Vars["API_KEY"] = "secret"
	base.Vars["UNIQUE"] = "only-here"

	ref1, _ := New("app", "staging")
	ref1.Vars["DB_HOST"] = "staging-host"

	ref2, _ := New("app", "production")
	ref2.Vars["API_KEY"] = "prod-secret"

	return base, ref1, ref2
}

func TestDedup_KeepFirst(t *testing.T) {
	base, ref1, ref2 := baseDedupSets(t)
	res, err := Dedup(base, []*EnvSet{ref1, ref2}, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if base.Vars["DB_HOST"] != "localhost" {
		t.Errorf("keep-first: expected original DB_HOST value")
	}
	if base.Vars["API_KEY"] != "secret" {
		t.Errorf("keep-first: expected original API_KEY value")
	}
}

func TestDedup_KeepLast(t *testing.T) {
	base, ref1, ref2 := baseDedupSets(t)
	_, err := Dedup(base, []*EnvSet{ref1, ref2}, DedupKeepLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base.Vars["DB_HOST"] != "staging-host" {
		t.Errorf("keep-last: expected ref DB_HOST value, got %q", base.Vars["DB_HOST"])
	}
	if base.Vars["API_KEY"] != "prod-secret" {
		t.Errorf("keep-last: expected ref API_KEY value, got %q", base.Vars["API_KEY"])
	}
}

func TestDedup_UniqueKeyUntouched(t *testing.T) {
	base, ref1, ref2 := baseDedupSets(t)
	_, err := Dedup(base, []*EnvSet{ref1, ref2}, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base.Vars["UNIQUE"] != "only-here" {
		t.Errorf("unique key should be untouched")
	}
}

func TestDedup_NilBase(t *testing.T) {
	_, err := Dedup(nil, nil, DedupKeepFirst)
	if err == nil {
		t.Error("expected error for nil base")
	}
}

func TestDedup_UnknownStrategy(t *testing.T) {
	base, _, _ := baseDedupSets(t)
	_, err := Dedup(base, nil, DedupStrategy("unknown"))
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestDedup_NilRefSkipped(t *testing.T) {
	base, _, _ := baseDedupSets(t)
	res, err := Dedup(base, []*EnvSet{nil}, DedupKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("nil ref should be skipped, got %d removed", len(res.Removed))
	}
}
