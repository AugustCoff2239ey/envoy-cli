package envset

import "fmt"

// RenameOptions configures the rename operation.
type RenameOptions struct {
	// NewName sets a new name for the envset. Leave empty to keep current.
	NewName string
	// NewEnvironment sets a new environment. Leave empty to keep current.
	NewEnvironment string
}

// Rename creates a copy of src with updated name/environment metadata,
// removes the old entry from the store, and saves the new one.
// At least one of opts.NewName or opts.NewEnvironment must differ from src.
func Rename(store *Store, src *EnvSet, opts RenameOptions) (*EnvSet, error) {
	if src == nil {
		return nil, fmt.Errorf("rename: source envset must not be nil")
	}

	newName := src.Name
	if opts.NewName != "" {
		newName = opts.NewName
	}

	newEnv := src.Environment
	if opts.NewEnvironment != "" {
		newEnv = opts.NewEnvironment
	}

	if newName == src.Name && newEnv == src.Environment {
		return nil, fmt.Errorf("rename: new name and environment are identical to current")
	}

	// Validate the new identity.
	candidate, err := New(newName, newEnv)
	if err != nil {
		return nil, fmt.Errorf("rename: %w", err)
	}

	// Copy all vars.
	for k, v := range src.Vars {
		candidate.Vars[k] = v
	}

	// Persist new, remove old.
	if err := store.Save(candidate); err != nil {
		return nil, fmt.Errorf("rename: save failed: %w", err)
	}

	if err := store.Delete(src.Name, src.Environment); err != nil {
		// Best-effort rollback.
		_ = store.Delete(candidate.Name, candidate.Environment)
		return nil, fmt.Errorf("rename: delete old failed: %w", err)
	}

	return candidate, nil
}
