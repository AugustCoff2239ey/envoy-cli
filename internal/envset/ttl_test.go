package envset

import (
	"testing"
	"time"
)

func baseTTLSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("ttlset", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["API_KEY"] = "abc123"
	es.Vars["DB_PASS"] = "secret"
	return es
}

func TestSetTTL_Valid(t *testing.T) {
	es := baseTTLSet(t)
	if err := SetTTL(es, "API_KEY", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ttl, err := GetTTL(es, "API_KEY")
	if err != nil {
		t.Fatalf("GetTTL: %v", err)
	}
	if ttl <= 0 || ttl > 5*time.Minute {
		t.Errorf("unexpected TTL: %v", ttl)
	}
}

func TestSetTTL_NilEnvSet(t *testing.T) {
	if err := SetTTL(nil, "API_KEY", time.Minute); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestSetTTL_MissingKey(t *testing.T) {
	es := baseTTLSet(t)
	if err := SetTTL(es, "MISSING", time.Minute); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSetTTL_ZeroDuration(t *testing.T) {
	es := baseTTLSet(t)
	if err := SetTTL(es, "API_KEY", 0); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestGetTTL_Expired(t *testing.T) {
	es := baseTTLSet(t)
	if err := SetExpiry(es, "DB_PASS", time.Now().Add(-time.Second)); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}
	if _, err := GetTTL(es, "DB_PASS"); err == nil {
		t.Error("expected error for expired key")
	}
}

func TestPurgeTTL_RemovesExpired(t *testing.T) {
	es := baseTTLSet(t)
	_ = SetExpiry(es, "API_KEY", time.Now().Add(-time.Second))
	removed, err := PurgeTTL(es)
	if err != nil {
		t.Fatalf("PurgeTTL: %v", err)
	}
	if len(removed) != 1 || removed[0] != "API_KEY" {
		t.Errorf("expected [API_KEY] removed, got %v", removed)
	}
	if _, ok := es.Vars["API_KEY"]; ok {
		t.Error("API_KEY should have been purged")
	}
}
