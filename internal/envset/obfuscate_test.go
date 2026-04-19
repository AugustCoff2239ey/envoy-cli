package envset

import (
	"strings"
	"testing"
)

func baseObfuscateSet() *EnvSet {
	es, _ := New("obfuscate-test", "staging")
	_ = es.Set("API_KEY", "supersecret")
	_ = es.Set("DB_PASS", "hunter2")
	_ = es.Set("APP_NAME", "envoy")
	return es
}

func TestObfuscate_AllKeys_Asterisks(t *testing.T) {
	es := baseObfuscateSet()
	out, err := Obfuscate(es, DefaultObfuscateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out.Vars {
		if !strings.Contains(v, "*") && v != "" {
			t.Errorf("key %s: expected asterisks, got %q", k, v)
		}
	}
}

func TestObfuscate_SelectedKeys(t *testing.T) {
	es := baseObfuscateSet()
	opts := ObfuscateOptions{Keys: []string{"API_KEY"}}
	out, err := Obfuscate(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.Vars["API_KEY"], "*") {
		t.Errorf("API_KEY should be obfuscated")
	}
	if out.Vars["APP_NAME"] != "envoy" {
		t.Errorf("APP_NAME should be unchanged, got %q", out.Vars["APP_NAME"])
	}
}

func TestObfuscate_UseHash(t *testing.T) {
	es := baseObfuscateSet()
	opts := ObfuscateOptions{UseHash: true}
	out, err := Obfuscate(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range out.Vars {
		if !strings.HasPrefix(v, "sha256:") {
			t.Errorf("expected sha256 prefix, got %q", v)
		}
	}
}

func TestObfuscate_EmptyValue(t *testing.T) {
	es, _ := New("empty-val", "local")
	_ = es.Set("EMPTY", "")
	out, err := Obfuscate(es, DefaultObfuscateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Vars["EMPTY"] != "" {
		t.Errorf("empty value should remain empty")
	}
}

func TestObfuscate_NilEnvSet(t *testing.T) {
	_, err := Obfuscate(nil, DefaultObfuscateOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestObfuscate_OriginalUnchanged(t *testing.T) {
	es := baseObfuscateSet()
	original := es.Vars["API_KEY"]
	_, _ = Obfuscate(es, DefaultObfuscateOptions())
	if es.Vars["API_KEY"] != original {
		t.Error("original EnvSet should not be mutated")
	}
}
