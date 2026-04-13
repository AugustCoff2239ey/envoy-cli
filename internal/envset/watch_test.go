package envset

import (
	"testing"
	"time"
)

func baseWatchSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("watch-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("HOST", "localhost")
	_ = es.Set("PORT", "8080")
	return es
}

func TestWatch_DetectsAdded(t *testing.T) {
	es := baseWatchSet(t)
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080"}

	stop := make(chan struct{})
	opts := WatchOptions{PollInterval: 20 * time.Millisecond, MaxEvents: 1}

	_ = es.Set("NEW_KEY", "new_value")

	ch, err := Watch(es, baseline, opts, stop)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.EventType != "added" {
			t.Errorf("expected added, got %s", ev.EventType)
		}
		if ev.Key != "NEW_KEY" {
			t.Errorf("expected NEW_KEY, got %s", ev.Key)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
	close(stop)
}

func TestWatch_DetectsChanged(t *testing.T) {
	es := baseWatchSet(t)
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080"}

	_ = es.Set("PORT", "9090")

	stop := make(chan struct{})
	opts := WatchOptions{PollInterval: 20 * time.Millisecond, MaxEvents: 1}

	ch, err := Watch(es, baseline, opts, stop)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.EventType != "changed" {
			t.Errorf("expected changed, got %s", ev.EventType)
		}
		if ev.OldValue != "8080" || ev.NewValue != "9090" {
			t.Errorf("unexpected values: old=%s new=%s", ev.OldValue, ev.NewValue)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
	close(stop)
}

func TestWatch_DetectsRemoved(t *testing.T) {
	es := baseWatchSet(t)
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080", "GONE": "bye"}

	stop := make(chan struct{})
	opts := WatchOptions{PollInterval: 20 * time.Millisecond, MaxEvents: 1}

	ch, err := Watch(es, baseline, opts, stop)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.EventType != "removed" {
			t.Errorf("expected removed, got %s", ev.EventType)
		}
		if ev.Key != "GONE" {
			t.Errorf("expected GONE, got %s", ev.Key)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
	close(stop)
}

func TestWatch_NilEnvSet(t *testing.T) {
	stop := make(chan struct{})
	defer close(stop)
	_, err := Watch(nil, nil, DefaultWatchOptions(), stop)
	if err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestWatch_StopSignal(t *testing.T) {
	es := baseWatchSet(t)
	baseline := map[string]string{"HOST": "localhost", "PORT": "8080"}

	stop := make(chan struct{})
	opts := WatchOptions{PollInterval: 20 * time.Millisecond, MaxEvents: 0}

	ch, err := Watch(es, baseline, opts, stop)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	close(stop)

	// Channel should close without hanging
	select {
	case _, ok := <-ch:
		_ = ok // channel closed cleanly
	case <-time.After(300 * time.Millisecond):
		t.Fatal("channel did not close after stop signal")
	}
}
