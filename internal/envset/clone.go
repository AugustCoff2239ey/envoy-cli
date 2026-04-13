package envset

import "fmt"

// CloneOptions configures the behavior of a Clone operation.
type CloneOptions struct {
	// NewName is the name for the cloned EnvSet.
	NewName string
	// NewEnvironment is the environment for the cloned EnvSet.
	NewEnvironment string
	// OverwriteExisting allows overwriting an existing EnvSet with the same name/env.
	OverwriteExisting bool
}

// Clone creates a deep copy of src using the provided options.
// The cloned EnvSet has its own independent copy of all key-value pairs.
func Clone(src *EnvSet, opts CloneOptions) (*EnvSet, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source EnvSet must not be nil")
	}

	name := opts.NewName
	if name == "" {
		name = src.Name
	}

	env := opts.NewEnvironment
	if env == "" {
		env = src.Environment
	}

	dst, err := New(name, env)
	if err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}

	for k, v := range src.Values {
		if err := dst.Set(k, v); err != nil {
			return nil, fmt.Errorf("clone: copying key %q: %w", k, err)
		}
	}

	return dst, nil
}
