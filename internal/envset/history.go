package envset

import (
	"errors"
	"time"
)

// HistoryEntry represents a single recorded change to an EnvSet.
type HistoryEntry struct {
	Timestamp time.Time
	Action    string
	Key       string
	OldValue  string
	NewValue  string
}

// History holds an ordered log of changes made to an EnvSet.
type History struct {
	entries []HistoryEntry
}

// NewHistory creates and returns an empty History.
func NewHistory() *History {
	return &History{
		entries: []HistoryEntry{},
	}
}

// Record appends a new entry to the history log.
func (h *History) Record(action, key, oldValue, newValue string) error {
	if h == nil {
		return errors.New("history: nil receiver")
	}
	if action == "" {
		return errors.New("history: action must not be empty")
	}
	h.entries = append(h.entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldValue,
		NewValue:  newValue,
	})
	return nil
}

// Entries returns a copy of all recorded history entries.
func (h *History) Entries() []HistoryEntry {
	if h == nil {
		return nil
	}
	copy := make([]HistoryEntry, len(h.entries))
	for i, e := range h.entries {
		copy[i] = e
	}
	return copy
}

// FilterByKey returns all entries matching the given key.
func (h *History) FilterByKey(key string) []HistoryEntry {
	var result []HistoryEntry
	for _, e := range h.entries {
		if e.Key == key {
			result = append(result, e)
		}
	}
	return result
}

// FilterByAction returns all entries matching the given action.
func (h *History) FilterByAction(action string) []HistoryEntry {
	var result []HistoryEntry
	for _, e := range h.entries {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// Clear removes all entries from the history.
func (h *History) Clear() {
	if h != nil {
		h.entries = []HistoryEntry{}
	}
}
