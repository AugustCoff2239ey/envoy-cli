package envset

import (
	"fmt"
	"time"
)

// Baseline represents a reference snapshot of an EnvSet used for drift detection.
type Baseline struct {
	Label     string
	CreatedAt time.Time
	Vars      map[string]string
}

// SetBaseline captures the current state of an EnvSet as a named baseline.
func SetBaseline(e *EnvSet, label string) (*Baseline, error) {
	if e == nil {
		return nil, fmt.Errorf("envset is nil")
	}
	if label == "" {
		return nil, fmt.Errorf("baseline label must not be empty")
	}

	vars := make(map[string]string, len(e.Vars))
	for k, v := range e.Vars {
		vars[k] = v
	}

	return &Baseline{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Vars:      vars,
	}, nil
}

// DriftResult describes how an EnvSet has changed relative to a Baseline.
type DriftResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [baseline value, current value]
}

// HasDrift returns true if any differences exist.
func (d *DriftResult) HasDrift() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// DetectDrift compares an EnvSet against a previously captured Baseline.
func DetectDrift(e *EnvSet, b *Baseline) (*DriftResult, error) {
	if e == nil {
		return nil, fmt.Errorf("envset is nil")
	}
	if b == nil {
		return nil, fmt.Errorf("baseline is nil")
	}

	result := &DriftResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, bv := range b.Vars {
		cv, ok := e.Vars[k]
		if !ok {
			result.Removed[k] = bv
		} else if cv != bv {
			result.Changed[k] = [2]string{bv, cv}
		}
	}

	for k, cv := range e.Vars {
		if _, ok := b.Vars[k]; !ok {
			result.Added[k] = cv
		}
	}

	return result, nil
}
