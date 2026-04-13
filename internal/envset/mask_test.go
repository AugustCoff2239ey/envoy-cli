package envset

import (
	"testing"
)

func baseMaskSet() *EnvSet {
	es, _ := New("myapp", "production")
	es.Vars["DB_PASSWORD"] = "supersecret1234"
	es.Vars["API_KEY"] = "abcdefgh"
	es.Vars["APP_ENV"] = "production"
	es.Vars["PORT"] = "8080"
	return es
}

func TestMaskValue_LongValue(t *testing.T) {
	result := MaskValue("supersecret1234")
	expected := "***********1234"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestMaskValue_ShortValue(t *testing.T) {
	result := MaskValue("abc")
	if result != "***" {
		t.Errorf("expected ***, got %q", result)
	}
}

func TestMaskValue_ExactlyFour(t *testing.T) {
	result := MaskValue("1234")
	if result != "****" {
		t.Errorf("expected ****, got %q", result)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "SECRET", "PRIVATE_KEY"}
	for _, k := range sensitive {
		if !IsSensitiveKey(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}

	notSensitive := []string{"PORT", "APP_ENV", "HOST", "DEBUG"}
	for _, k := range notSensitive {
		if IsSensitiveKey(k) {
			t.Errorf("expected %q to not be sensitive", k)
		}
	}
}

func TestMaskSensitive_MasksCorrectKeys(t *testing.T) {
	es := baseMaskSet()
	masked, err := MaskSensitive(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if masked.Vars["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unmasked")
	}
	if masked.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT to be unmasked")
	}
	if masked.Vars["DB_PASSWORD"] == "supersecret1234" {
		t.Errorf("expected DB_PASSWORD to be masked")
	}
	if masked.Vars["API_KEY"] == "abcdefgh" {
		t.Errorf("expected API_KEY to be masked")
	}
}

func TestMaskSensitive_NilSource(t *testing.T) {
	_, err := MaskSensitive(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestMaskKeys_SpecificKeys(t *testing.T) {
	es := baseMaskSet()
	masked, err := MaskKeys(es, []string{"PORT", "APP_ENV"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if masked.Vars["PORT"] == "8080" {
		t.Errorf("expected PORT to be masked")
	}
	if masked.Vars["APP_ENV"] == "production" {
		t.Errorf("expected APP_ENV to be masked")
	}
	if masked.Vars["DB_PASSWORD"] != "supersecret1234" {
		t.Errorf("expected DB_PASSWORD to remain unmasked")
	}
}

func TestMaskKeys_NilSource(t *testing.T) {
	_, err := MaskKeys(nil, []string{"PORT"})
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
