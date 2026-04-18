package envset

import (
	"strings"
	"testing"
)

func baseQuotaSet() *EnvSet {
	es, _ := New("quota-test", "local")
	_ = es.Set("KEY_ONE", "value1")
	_ = es.Set("KEY_TWO", "value2")
	_ = es.Set("KEY_THREE", "value3")
	return es
}

func TestCheckQuota_WithinLimits(t *testing.T) {
	es := baseQuotaSet()
	violations, err := CheckQuota(es, DefaultQuotaOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheckQuota_TooManyKeys(t *testing.T) {
	es := baseQuotaSet()
	opts := DefaultQuotaOptions()
	opts.MaxKeys = 2
	violations, err := CheckQuota(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 || violations[0].Key != "__count__" {
		t.Errorf("expected count violation, got %+v", violations)
	}
}

func TestCheckQuota_KeyTooLong(t *testing.T) {
	es, _ := New("quota-test", "local")
	longKey := strings.Repeat("A", 10)
	_ = es.Set(longKey, "val")
	opts := DefaultQuotaOptions()
	opts.MaxKeyLen = 5
	violations, err := CheckQuota(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheckQuota_ValueTooLong(t *testing.T) {
	es, _ := New("quota-test", "local")
	_ = es.Set("KEY", strings.Repeat("x", 20))
	opts := DefaultQuotaOptions()
	opts.MaxValLen = 10
	violations, err := CheckQuota(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheckQuota_NilEnvSet(t *testing.T) {
	_, err := CheckQuota(nil, DefaultQuotaOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
