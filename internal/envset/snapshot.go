package envset

import (
	"fmt"
	"time"
)

// Snapshot represents a point-in-time capture of an EnvSet's variables.
type Snapshot struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Env       string            `json:"env"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
	Message   string            `json:"message,omitempty"`
}

// TakeSnapshot creates a snapshot of the given EnvSet with an optional message.
func TakeSnapshot(es *EnvSet, message string) (*Snapshot, error) {
	if es == nil {
		return nil, fmt.Errorf("snapshot: source EnvSet must not be nil")
	}

	vars := make(map[string]string, len(es.Vars))
	for k, v := range es.Vars {
		vars[k] = v
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("%s-%s-%d", es.Name, es.Environment, now.UnixNano())

	return &Snapshot{
		ID:        id,
		Name:      es.Name,
		Env:       es.Environment,
		Vars:      vars,
		CreatedAt: now,
		Message:   message,
	}, nil
}

// RestoreSnapshot returns a new EnvSet populated from the given snapshot.
func RestoreSnapshot(snap *Snapshot) (*EnvSet, error) {
	if snap == nil {
		return nil, fmt.Errorf("snapshot: cannot restore from nil snapshot")
	}

	es, err := New(snap.Name, snap.Env)
	if err != nil {
		return nil, fmt.Errorf("snapshot: restore failed: %w", err)
	}

	for k, v := range snap.Vars {
		if err := es.Set(k, v); err != nil {
			return nil, fmt.Errorf("snapshot: restore failed setting %q: %w", k, err)
		}
	}

	return es, nil
}
