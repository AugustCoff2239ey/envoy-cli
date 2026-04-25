package envset

import "fmt"

// RenameKeyOptions controls the behaviour of RenameKey.
type RenameKeyOptions struct {
	// Overwrite allows the target key name to replace an existing key.
	Overwrite bool
}

// DefaultRenameKeyOptions returns safe defaults.
func DefaultRenameKeyOptions() RenameKeyOptions {
	return RenameKeyOptions{Overwrite: false}
}

// RenameKey renames a key within an EnvSet from oldKey to newKey.
// The value, metadata (lock, pin, annotations, labels) are preserved.
// Returns an error if:
//   - the EnvSet is nil or frozen / readonly
//   - oldKey does not exist
//   - newKey is invalid
//   - newKey already exists and Overwrite is false
func RenameKey(es *EnvSet, oldKey, newKey string, opts RenameKeyOptions) error {
	if es == nil {
		return fmt.Errorf("rename_key: nil EnvSet")
	}
	if err := AssertMutable(es); err != nil {
		return fmt.Errorf("rename_key: %w", err)
	}
	if err := AssertWritable(es); err != nil {
		return fmt.Errorf("rename_key: %w", err)
	}
	if err := ValidateKey(oldKey); err != nil {
		return fmt.Errorf("rename_key: invalid old key: %w", err)
	}
	if err := ValidateKey(newKey); err != nil {
		return fmt.Errorf("rename_key: invalid new key: %w", err)
	}

	val, ok := es.Vars[oldKey]
	if !ok {
		return fmt.Errorf("rename_key: key %q not found", oldKey)
	}
	if IsLocked(es, oldKey) {
		return fmt.Errorf("rename_key: key %q is locked", oldKey)
	}
	if IsProtected(es, oldKey) {
		return fmt.Errorf("rename_key: key %q is protected", oldKey)
	}

	if _, exists := es.Vars[newKey]; exists && !opts.Overwrite {
		return fmt.Errorf("rename_key: key %q already exists (use overwrite to replace)", newKey)
	}

	// Move value.
	es.Vars[newKey] = val
	delete(es.Vars, oldKey)

	// Transfer metadata stored under the old key name.
	transferMeta(es, oldKey, newKey)

	return nil
}

// transferMeta copies key-level metadata from src to dst and removes src entries.
func transferMeta(es *EnvSet, src, dst string) {
	if es.Meta == nil {
		return
	}
	if m, ok := es.Meta[src]; ok {
		es.Meta[dst] = m
		delete(es.Meta, src)
	}
}
