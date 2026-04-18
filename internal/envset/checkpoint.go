package envset

import (
	"errors"
	"fmt"
	"time"
)

// Checkpoint represents a named point-in-time snapshot of an EnvSet.
type Checkpoint struct {
	Name      string
	CreatedAt time.Time
	Vars      map[string]string
}

// CheckpointStore holds named checkpoints for an EnvSet.
type CheckpointStore struct {
	checkpoints map[string]*Checkpoint
}

// NewCheckpointStore initialises an empty CheckpointStore.
func NewCheckpointStore() *CheckpointStore {
	return &CheckpointStore{checkpoints: make(map[string]*Checkpoint)}
}

// Save creates or overwrites a named checkpoint from the given EnvSet.
func (cs *CheckpointStore) Save(name string, e *EnvSet) error {
	if e == nil {
		return errors.New("checkpoint: nil EnvSet")
	}
	if name == "" {
		return errors.New("checkpoint: name must not be empty")
	}
	vars := make(map[string]string, len(e.Vars))
	for k, v := range e.Vars {
		vars[k] = v
	}
	cs.checkpoints[name] = &Checkpoint{
		Name:      name,
		CreatedAt: time.Now(),
		Vars:      vars,
	}
	return nil
}

// Load retrieves a named checkpoint.
func (cs *CheckpointStore) Load(name string) (*Checkpoint, error) {
	cp, ok := cs.checkpoints[name]
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", name)
	}
	return cp, nil
}

// Restore applies a named checkpoint's vars onto the given EnvSet.
func (cs *CheckpointStore) Restore(name string, e *EnvSet) error {
	if e == nil {
		return errors.New("checkpoint: nil EnvSet")
	}
	cp, err := cs.Load(name)
	if err != nil {
		return err
	}
	e.Vars = make(map[string]string, len(cp.Vars))
	for k, v := range cp.Vars {
		e.Vars[k] = v
	}
	return nil
}

// List returns all checkpoint names in the store.
func (cs *CheckpointStore) List() []string {
	names := make([]string, 0, len(cs.checkpoints))
	for n := range cs.checkpoints {
		names = append(names, n)
	}
	return names
}

// Delete removes a named checkpoint.
func (cs *CheckpointStore) Delete(name string) error {
	if _, ok := cs.checkpoints[name]; !ok {
		return fmt.Errorf("checkpoint %q not found", name)
	}
	delete(cs.checkpoints, name)
	return nil
}
