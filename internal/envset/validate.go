package envset

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// keyPattern matches valid environment variable key names.
	// Keys must start with a letter or underscore, followed by letters, digits, or underscores.
	keyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

	// reservedKeys are keys that should not be set directly.
	reservedKeys = map[string]bool{
		"PATH": true,
		"HOME": true,
		"USER": true,
	}
)

// ValidationError holds a list of validation issues found in an EnvSet.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d issue(s): %s", len(e.Issues), strings.Join(e.Issues, "; "))
}

// Validate checks all keys and values in the given EnvSet for common issues.
// It returns a *ValidationError if any issues are found, or nil if the set is valid.
func Validate(es *EnvSet) error {
	var issues []string

	for key, value := range es.Vars {
		if !keyPattern.MatchString(key) {
			issues = append(issues, fmt.Sprintf("invalid key name %q: must match [A-Za-z_][A-Za-z0-9_]*", key))
		}

		if reservedKeys[strings.ToUpper(key)] {
			issues = append(issues, fmt.Sprintf("key %q is reserved and should not be managed by envoy", key))
		}

		if strings.Contains(value, "\n") {
			issues = append(issues, fmt.Sprintf("value for key %q contains a newline character", key))
		}

		if len(value) > 4096 {
			issues = append(issues, fmt.Sprintf("value for key %q exceeds maximum length of 4096 characters", key))
		}
	}

	if len(issues) == 0 {
		return nil
	}
	return &ValidationError{Issues: issues}
}

// ValidateKey returns an error if the provided key is not a valid environment variable name.
func ValidateKey(key string) error {
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("invalid key name %q: must match [A-Za-z_][A-Za-z0-9_]*", key)
	}
	return nil
}
