package envset

import (
	"errors"
	"regexp"
)

// Environment represents a target deployment environment.
type Environment string

const (
	EnvLocal      Environment = "local"
	EnvStaging    Environment = "staging"
	EnvProduction Environment = "production"
)

var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// EnvSet holds a named collection of environment variables for a specific environment.
type EnvSet struct {
	Name        string            `json:"name"`
	Environment Environment       `json:"environment"`
	Variables   map[string]string `json:"variables"`
}

// New creates a new EnvSet with the given name and environment.
func New(name string, env Environment) (*EnvSet, error) {
	if name == "" {
		return nil, errors.New("envset name must not be empty")
	}
	if !isValidEnvironment(env) {
		return nil, errors.New("invalid environment: must be local, staging, or production")
	}
	return &EnvSet{
		Name:        name,
		Environment: env,
		Variables:   make(map[string]string),
	}, nil
}

// Set adds or updates a key-value pair in the EnvSet.
func (e *EnvSet) Set(key, value string) error {
	if !validKeyPattern.MatchString(key) {
		return errors.New("invalid key: must be uppercase letters, digits, or underscores, and not start with a digit")
	}
	e.Variables[key] = value
	return nil
}

// Get retrieves the value for a given key.
func (e *EnvSet) Get(key string) (string, bool) {
	v, ok := e.Variables[key]
	return v, ok
}

// Delete removes a key from the EnvSet.
func (e *EnvSet) Delete(key string) {
	delete(e.Variables, key)
}

// Keys returns all keys in the EnvSet.
func (e *EnvSet) Keys() []string {
	keys := make([]string, 0, len(e.Variables))
	for k := range e.Variables {
		keys = append(keys, k)
	}
	return keys
}

func isValidEnvironment(env Environment) bool {
	switch env {
	case EnvLocal, EnvStaging, EnvProduction:
		return true
	}
	return false
}
