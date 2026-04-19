package envset

import (
	"testing"
	"time"
)

func baseSummarySet(t *testing.T) *EnvSet {
	t.Helper()
	es, _ := New("summary-test", "staging")
	_ = es.Set("API_KEY", "secret")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("PORT", "8080")
	return es
}

func TestSummary_Basic(t *testing.T) {
	es := baseSummarySet(t)
	r, err := Summary(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 keys, got %d", r.TotalKeys)
	}
	if r.Name != "summary-test" || r.Environment != "staging" {
		t.Errorf("unexpected name/env: %s/%s", r.Name, r.Environment)
	}
}

func TestSummary_LockedAndProtected(t *testing.T) {
	es := baseSummarySet(t)
	_ = LockKey(es, "API_KEY", "test")
	_ = ProtectKey(es, "DB_HOST")

	r, err := Summary(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.LockedKeys) != 1 || r.LockedKeys[0] != "API_KEY" {
		t.Errorf("expected API_KEY locked, got %v", r.LockedKeys)
	}
	if len(r.Protected) != 1 || r.Protected[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST protected, got %v", r.Protected)
	}
}

func TestSummary_ReadonlyFrozen(t *testing.T) {
	es := baseSummarySet(t)
	MarkReadonly(es)
	Freeze(es)

	r, err := Summary(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Readonly {
		t.Error("expected readonly=true")
	}
	if !r.Frozen {
		t.Error("expected frozen=true")
	}
}

func TestSummary_Expired(t *testing.T) {
	es := baseSummarySet(t)
	_ = SetExpiry(es, "PORT", time.Now().Add(-time.Hour))

	r, err := Summary(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Expired) != 1 || r.Expired[0] != "PORT" {
		t.Errorf("expected PORT expired, got %v", r.Expired)
	}
}

func TestSummary_NilEnvSet(t *testing.T) {
	_, err := Summary(nil)
	if err == nil {
		t.Error("expected error for nil envset")
	}
}
