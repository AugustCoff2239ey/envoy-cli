package envset

import (
	"testing"
)

func validatableSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("myapp", "staging")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	es.Vars["DB_HOST"] = "localhost"
	es.Vars["API_KEY"] = "abc123"
	return es
}

func TestValidate_ValidSet(t *testing.T) {
	es := validatableSet(t)
	if err := Validate(es); err != nil {
		t.Errorf("Validate() expected nil error, got: %v", err)
	}
}

func TestValidate_InvalidKey(t *testing.T) {
	es := validatableSet(t)
	es.Vars["123INVALID"] = "value"

	err := Validate(es)
	if err == nil {
		t.Fatal("Validate() expected error for invalid key, got nil")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) == 0 {
		t.Error("expected at least one issue")
	}
}

func TestValidate_ReservedKey(t *testing.T) {
	es := validatableSet(t)
	es.Vars["PATH"] = "/usr/local/bin"

	err := Validate(es)
	if err == nil {
		t.Fatal("Validate() expected error for reserved key, got nil")
	}
}

func TestValidate_NewlineInValue(t *testing.T) {
	es := validatableSet(t)
	es.Vars["MULTILINE"] = "line1\nline2"

	err := Validate(es)
	if err == nil {
		t.Fatal("Validate() expected error for newline in value, got nil")
	}
}

func TestValidate_ValueTooLong(t *testing.T) {
	es := validatableSet(t)
	long := make([]byte, 4097)
	for i := range long {
		long[i] = 'x'
	}
	es.Vars["BIG_VALUE"] = string(long)

	err := Validate(es)
	if err == nil {
		t.Fatal("Validate() expected error for oversized value, got nil")
	}
}

func TestValidateKey_Valid(t *testing.T) {
	cases := []string{"MY_VAR", "_PRIVATE", "A", "VAR123"}
	for _, key := range cases {
		if err := ValidateKey(key); err != nil {
			t.Errorf("ValidateKey(%q) unexpected error: %v", key, err)
		}
	}
}

func TestValidateKey_Invalid(t *testing.T) {
	cases := []string{"1INVALID", "has-dash", "has space", ""}
	for _, key := range cases {
		if err := ValidateKey(key); err == nil {
			t.Errorf("ValidateKey(%q) expected error, got nil", key)
		}
	}
}
