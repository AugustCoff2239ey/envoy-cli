package envset

import (
	"testing"
)

func baseSchemaSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("schema-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["PORT"] = "8080"
	es.Vars["DEBUG"] = "true"
	es.Vars["API_URL"] = "https://api.example.com"
	es.Vars["APP_NAME"] = "myapp"
	return es
}

func TestValidateSchema_ValidSet(t *testing.T) {
	es := baseSchemaSet(t)
	schema := Schema{
		"PORT":     {Type: TypeInteger, Required: true},
		"DEBUG":    {Type: TypeBoolean, Required: true},
		"API_URL":  {Type: TypeURL, Required: true},
		"APP_NAME": {Type: TypeString, Required: true},
	}
	violations, err := ValidateSchema(es, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	es := baseSchemaSet(t)
	schema := Schema{
		"MISSING_KEY": {Type: TypeString, Required: true},
	}
	violations, err := ValidateSchema(es, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 || violations[0].Key != "MISSING_KEY" {
		t.Errorf("expected violation for MISSING_KEY, got %v", violations)
	}
}

func TestValidateSchema_WrongType(t *testing.T) {
	es := baseSchemaSet(t)
	es.Vars["PORT"] = "not-a-number"
	schema := Schema{
		"PORT": {Type: TypeInteger, Required: true},
	}
	violations, err := ValidateSchema(es, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Error("expected type violation for PORT")
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	es := baseSchemaSet(t)
	es.Vars["APP_NAME"] = "My App!"
	schema := Schema{
		"APP_NAME": {Type: TypeString, Pattern: `^[a-z0-9_-]+$`},
	}
	violations, err := ValidateSchema(es, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) == 0 {
		t.Error("expected pattern violation for APP_NAME")
	}
}

func TestValidateSchema_OptionalMissing(t *testing.T) {
	es := baseSchemaSet(t)
	schema := Schema{
		"OPTIONAL_KEY": {Type: TypeString, Required: false},
	}
	violations, err := ValidateSchema(es, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations for optional missing key, got %v", violations)
	}
}

func TestValidateSchema_NilEnvSet(t *testing.T) {
	_, err := ValidateSchema(nil, Schema{})
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestValidateSchema_InvalidPattern(t *testing.T) {
	es := baseSchemaSet(t)
	schema := Schema{
		"APP_NAME": {Type: TypeString, Pattern: `[invalid(`},
	}
	_, err := ValidateSchema(es, schema)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestValidateSchema_EmptySchema(t *testing.T) {
	es := baseSchemaSet(t)
	violations, err := ValidateSchema(es, Schema{})
	if err != nil {
		t.Fatalf("unexpected error for empty schema: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations for empty schema, got %v", violations)
	}
}
