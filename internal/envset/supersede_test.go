package envset

import (
	"testing"
)

func baseSupersedesets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	src, _ := New("source", "staging")
	src.Vars["DB_HOST"] = "prod-db.example.com"
	src.Vars["API_KEY"] = "new-secret"
	src.Vars["TIMEOUT"] = "60"

	dst, _ := New("target", "staging")
	dst.Vars["DB_HOST"] = "local-db"
	dst.Vars["API_KEY"] = "old-secret"
	dst.Vars["TIMEOUT"] = "30"
	return src, dst
}

func TestSupersede_AllKeys(t *testing.T) {
	src, dst := baseSupersedesets(t)
	opts := DefaultSupersedeOptions()
	n, err := Supersede(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 replaced, got %d", n)
	}
	if dst.Vars["DB_HOST"] != "prod-db.example.com" {
		t.Errorf("DB_HOST not superseded")
	}
}

func TestSupersede_SelectedKeys(t *testing.T) {
	src, dst := baseSupersedesets(t)
	opts := DefaultSupersedeOptions()
	opts.Keys = []string{"DB_HOST"}
	n, err := Supersede(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 replaced, got %d", n)
	}
	if dst.Vars["API_KEY"] != "old-secret" {
		t.Errorf("API_KEY should not have changed")
	}
}

func TestSupersede_SkipsLockedKey(t *testing.T) {
	src, dst := baseSupersedesets(t)
	_, _ = LockKey(dst, "API_KEY", "test")
	opts := DefaultSupersedeOptions()
	n, err := Supersede(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 replaced (locked skipped), got %d", n)
	}
	if dst.Vars["API_KEY"] != "old-secret" {
		t.Errorf("locked API_KEY should not have been superseded")
	}
}

func TestSupersede_SkipsProtectedKey(t *testing.T) {
	src, dst := baseSupersedesets(t)
	_ = ProtectKey(dst, "TIMEOUT")
	opts := DefaultSupersedeOptions()
	n, err := Supersede(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 replaced (protected skipped), got %d", n)
	}
	if dst.Vars["TIMEOUT"] != "30" {
		t.Errorf("protected TIMEOUT should not have been superseded")
	}
}

func TestSupersede_NilSource(t *testing.T) {
	_, dst := baseSupersedesets(t)
	_, err := Supersede(nil, dst, DefaultSupersedeOptions())
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestSupersede_MissingSourceKey(t *testing.T) {
	src, dst := baseSupersedesets(t)
	opts := DefaultSupersedeOptions()
	opts.Keys = []string{"NONEXISTENT"}
	_, err := Supersede(src, dst, opts)
	if err == nil {
		t.Error("expected error for missing source key")
	}
}
