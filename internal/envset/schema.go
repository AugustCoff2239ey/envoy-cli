package envset

import (
	"fmt"
	"regexp"
)

// SchemaFieldType represents the expected type of an env var value.
type SchemaFieldType string

const (
	TypeString  SchemaFieldType = "string"
	TypeInteger SchemaFieldType = "integer"
	TypeBoolean SchemaFieldType = "boolean"
	TypeURL     SchemaFieldType = "url"
)

// SchemaField defines constraints for a single environment variable.
type SchemaField struct {
	Type     SchemaFieldType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema maps env var keys to their field definitions.
type Schema map[string]SchemaField

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

var (
	reInteger = regexp.MustCompile(`^-?\d+$`)
	reBoolean = regexp.MustCompile(`^(true|false|1|0|yes|no)$`)
	reURL     = regexp.MustCompile(`^https?://[^\s]+$`)
)

// ValidateSchema checks an EnvSet against a Schema and returns any violations.
func ValidateSchema(es *EnvSet, schema Schema) ([]SchemaViolation, error) {
	if es == nil {
		return nil, fmt.Errorf("envset is nil")
	}

	var violations []SchemaViolation

	for key, field := range schema {
		val, exists := es.Vars[key]

		if field.Required && !exists {
			violations = append(violations, SchemaViolation{Key: key, Message: "required key is missing"})
			continue
		}

		if !exists {
			continue
		}

		switch field.Type {
		case TypeInteger:
			if !reInteger.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("expected integer, got %q", val)})
			}
		case TypeBoolean:
			if !reBoolean.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("expected boolean, got %q", val)})
			}
		case TypeURL:
			if !reURL.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("expected URL, got %q", val)})
			}
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for key %q: %w", key, err)
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern)})
			}
		}
	}

	return violations, nil
}
