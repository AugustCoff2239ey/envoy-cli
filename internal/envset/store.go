package envset

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const defaultStoreDir = ".envoy"

// Store manages persistence of EnvSets on disk.
type Store struct {
	BaseDir string
}

// NewStore creates a Store rooted at the given directory.
// If dir is empty, it defaults to ".envoy" in the current working directory.
func NewStore(dir string) (*Store, error) {
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("could not determine working directory: %w", err)
		}
		dir = filepath.Join(cwd, defaultStoreDir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("could not create store directory: %w", err)
	}
	return &Store{BaseDir: dir}, nil
}

// filePath returns the path for a given envset name and environment.
func (s *Store) filePath(name string, env Environment) string {
	filename := fmt.Sprintf("%s.%s.json", name, env)
	return filepath.Join(s.BaseDir, filename)
}

// Save persists an EnvSet to disk as JSON.
func (s *Store) Save(es *EnvSet) error {
	if es == nil {
		return errors.New("cannot save nil EnvSet")
	}
	data, err := json.MarshalIndent(es, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal envset: %w", err)
	}
	path := s.filePath(es.Name, es.Environment)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write envset file: %w", err)
	}
	return nil
}

// Load reads an EnvSet from disk by name and environment.
func (s *Store) Load(name string, env Environment) (*EnvSet, error) {
	path := s.filePath(name, env)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("envset '%s' for environment '%s' not found", name, env)
		}
		return nil, fmt.Errorf("failed to read envset file: %w", err)
	}
	var es EnvSet
	if err := json.Unmarshal(data, &es); err != nil {
		return nil, fmt.Errorf("failed to parse envset file: %w", err)
	}
	return &es, nil
}

// Delete removes an EnvSet file from disk.
func (s *Store) Delete(name string, env Environment) error {
	path := s.filePath(name, env)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("envset '%s' for environment '%s' does not exist", name, env)
		}
		return fmt.Errorf("failed to delete envset file: %w", err)
	}
	return nil
}
