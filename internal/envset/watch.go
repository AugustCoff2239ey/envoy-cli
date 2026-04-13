package envset

import (
	"fmt"
	"time"
)

// WatchEvent represents a change detected during a watch cycle.
type WatchEvent struct {
	Key       string
	OldValue  string
	NewValue  string
	EventType string // "added", "removed", "changed"
	At        time.Time
}

// WatchOptions configures the behavior of a Watch operation.
type WatchOptions struct {
	// PollInterval defines how often to check for changes.
	PollInterval time.Duration
	// MaxEvents stops watching after this many events (0 = unlimited).
	MaxEvents int
}

// DefaultWatchOptions returns sensible defaults for WatchOptions.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		PollInterval: 2 * time.Second,
		MaxEvents:    0,
	}
}

// Watch polls the given EnvSet against a baseline snapshot, emitting WatchEvents
// on the returned channel whenever changes are detected. The caller must close
// the stop channel to terminate the watch loop.
func Watch(current *EnvSet, baseline map[string]string, opts WatchOptions, stop <-chan struct{}) (<-chan WatchEvent, error) {
	if current == nil {
		return nil, fmt.Errorf("watch: current EnvSet must not be nil")
	}

	events := make(chan WatchEvent, 16)

	go func() {
		defer close(events)
		prev := copyVars(baseline)
		total := 0

		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				current.mu.RLock()
				now := copyVars(current.vars)
				current.mu.RUnlock()

				for k, newVal := range now {
					oldVal, existed := prev[k]
					if !existed {
						events <- WatchEvent{Key: k, OldValue: "", NewValue: newVal, EventType: "added", At: time.Now()}
						total++
					} else if oldVal != newVal {
						events <- WatchEvent{Key: k, OldValue: oldVal, NewValue: newVal, EventType: "changed", At: time.Now()}
						total++
					}
				}
				for k, oldVal := range prev {
					if _, ok := now[k]; !ok {
						events <- WatchEvent{Key: k, OldValue: oldVal, NewValue: "", EventType: "removed", At: time.Now()}
						total++
					}
				}

				prev = now

				if opts.MaxEvents > 0 && total >= opts.MaxEvents {
					return
				}
			}
		}
	}()

	return events, nil
}

// copyVars returns a shallow copy of a string map.
func copyVars(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
