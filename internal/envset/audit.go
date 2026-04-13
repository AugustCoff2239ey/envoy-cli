package envset

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change made to an EnvSet.
type AuditAction string

const (
	AuditActionSet    AuditAction = "set"
	AuditActionDelete AuditAction = "delete"
	AuditActionImport AuditAction = "import"
	AuditActionSync   AuditAction = "sync"
	AuditActionMerge  AuditAction = "merge"
)

// AuditEntry records a single change event on an EnvSet.
type AuditEntry struct {
	Timestamp   time.Time   `json:"timestamp"`
	Action      AuditAction `json:"action"`
	Environment string      `json:"environment"`
	Key         string      `json:"key,omitempty"`
	Description string      `json:"description"`
}

// AuditLog holds an ordered list of audit entries for an EnvSet.
type AuditLog struct {
	Name    string       `json:"name"`
	Entries []AuditEntry `json:"entries"`
}

// NewAuditLog creates an empty AuditLog for the given EnvSet name.
func NewAuditLog(name string) *AuditLog {
	return &AuditLog{
		Name:    name,
		Entries: []AuditEntry{},
	}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(action AuditAction, env, key, description string) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp:   time.Now().UTC(),
		Action:      action,
		Environment: env,
		Key:         key,
		Description: description,
	})
}

// Filter returns entries matching the given action. If action is empty, all entries are returned.
func (a *AuditLog) Filter(action AuditAction) []AuditEntry {
	if action == "" {
		return a.Entries
	}
	var out []AuditEntry
	for _, e := range a.Entries {
		if e.Action == action {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable summary of the audit log.
func (a *AuditLog) Summary() string {
	return fmt.Sprintf("AuditLog[%s]: %d entries", a.Name, len(a.Entries))
}
