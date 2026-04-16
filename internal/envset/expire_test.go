package envset

import (
	"testing"
	"time"
)

func baseExpireSet() *EnvSet {
	e, _ := New("expire-test", "local")
	e.Vars["API_KEY"] = "secret"
	e.Vars["DB_PASS"] = "pass123"
	return e
}

func TestSetExpiry_Valid(t *testing.T) {
	e := baseExpireSet()
	expiry := time.Now().Add(time.Hour)
	if err := SetExpiry(e, "API_KEY", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, set, err := GetExpiry(e, "API_KEY")
	if err != nil || !set {
		t.Fatalf("expected expiry to be set, err=%v", err)
	}
	if got.Unix() != expiry.UTC().Truncate(time.Second).Unix() {
		t.Errorf("expiry mismatch: got %v", got)
	}
}

func TestSetExpiry_NilEnvSet(t *testing.T) {
	if err := SetExpiry(nil, "API_KEY", time.Now()); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestSetExpiry_MissingKey(t *testing.T) {
	e := baseExpireSet()
	if err := SetExpiry(e, "MISSING", time.Now()); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestIsExpired_NotExpired(t *testing.T) {
	e := baseExpireSet()
	_ = SetExpiry(e, "API_KEY", time.Now().Add(time.Hour))
	expired, err := IsExpired(e, "API_KEY")
	if err != nil || expired {
		t.Errorf("expected not expired, err=%v", err)
	}
}

func TestIsExpired_Expired(t *testing.T) {
	e := baseExpireSet()
	_ = SetExpiry(e, "API_KEY", time.Now().Add(-time.Second))
	expired, err := IsExpired(e, "API_KEY")
	if err != nil || !expired {
		t.Errorf("expected expired, err=%v", err)
	}
}

func TestPurgeExpired_RemovesExpired(t *testing.T) {
	e := baseExpireSet()
	_ = SetExpiry(e, "API_KEY", time.Now().Add(-time.Second))
	_ = SetExpiry(e, "DB_PASS", time.Now().Add(time.Hour))
	purged, err := PurgeExpired(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(purged) != 1 || purged[0] != "API_KEY" {
		t.Errorf("expected API_KEY purged, got %v", purged)
	}
	if _, ok := e.Vars["API_KEY"]; ok {
		t.Error("expected API_KEY to be removed")
	}
	if _, ok := e.Vars["DB_PASS"]; !ok {
		t.Error("expected DB_PASS to remain")
	}
}

func TestPurgeExpired_NilEnvSet(t *testing.T) {
	_, err := PurgeExpired(nil)
	if err == nil {
		t.Error("expected error for nil envset")
	}
}
