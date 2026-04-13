package envset

import (
	"errors"
	"fmt"
)

// RollbackEntry represents a saved state that can be restored.
type RollbackEntry struct {
	Label   string
	Vars    map[string]string
	Message string
}

// RollbackStack holds a stack of rollback entries for an EnvSet.
type RollbackStack struct {
	entries []RollbackEntry
	max     int
}

// NewRollbackStack creates a new RollbackStack with a max capacity.
func NewRollbackStack(max int) *RollbackStack {
	if max <= 0 {
		max = 10
	}
	return &RollbackStack{max: max}
}

// Push saves the current state of an EnvSet onto the stack.
func (r *RollbackStack) Push(es *EnvSet, message string) error {
	if es == nil {
		return errors.New("rollback: envset is nil")
	}
	snapshot := make(map[string]string, len(es.Vars))
	for k, v := range es.Vars {
		snapshot[k] = v
	}
	entry := RollbackEntry{
		Label:   fmt.Sprintf("%s/%s", es.Name, es.Environment),
		Vars:    snapshot,
		Message: message,
	}
	if len(r.entries) >= r.max {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, entry)
	return nil
}

// Pop restores the most recent state into the given EnvSet.
func (r *RollbackStack) Pop(es *EnvSet) error {
	if es == nil {
		return errors.New("rollback: envset is nil")
	}
	if len(r.entries) == 0 {
		return errors.New("rollback: no entries to restore")
	}
	last := r.entries[len(r.entries)-1]
	r.entries = r.entries[:len(r.entries)-1]
	es.Vars = make(map[string]string, len(last.Vars))
	for k, v := range last.Vars {
		es.Vars[k] = v
	}
	return nil
}

// Len returns the number of entries in the stack.
func (r *RollbackStack) Len() int {
	return len(r.entries)
}

// Peek returns the most recent entry without removing it.
func (r *RollbackStack) Peek() (RollbackEntry, error) {
	if len(r.entries) == 0 {
		return RollbackEntry{}, errors.New("rollback: stack is empty")
	}
	return r.entries[len(r.entries)-1], nil
}
