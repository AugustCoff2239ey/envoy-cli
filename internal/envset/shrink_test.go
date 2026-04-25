package envset

import (
	"testing"
	"time"
)

func baseShrinkSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("shrink-set", "test")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, kv := range []struct{ k, v string }{
		{"ALPHA", "one"},
		{"BETA", "two"},
		{"GAMMA", ""},
		{"DELTA", "four"},
	} {
		if err := es.Set(kv.k, kv.v); err != nil {
			t.Fatalf("Set %q: %v", kv.k, err)
		}
	}
	return es
}

func TestShrink_RemoveEmpty(t *testing.T) {
	es := baseShrinkSet(t)
	opts := DefaultShrinkOptions()
	opts.RemoveExpired = false
	opts.RemoveEmpty = true

	if err := Shrink(es, opts); err != nil {
		t.Fatalf("Shrink: %v", err)
	}
	if _, err := es.Get("GAMMA"); err == nil {
		t.Error("expected GAMMA to be removed (empty value)")
	}
	if _, err := es.Get("ALPHA"); err != nil {
		t.Error("expected ALPHA to remain")
	}
}

func TestShrink_MaxKeys(t *testing.T) {
	es := baseShrinkSet(t)
	opts := DefaultShrinkOptions()
	opts.RemoveExpired = false
	opts.MaxKeys = 2

	if err := Shrink(es, opts); err != nil {
		t.Fatalf("Shrink: %v", err)
	}
	if got := len(es.Keys()); got != 2 {
		t.Errorf("expected 2 keys after shrink, got %d", got)
	}
}

func TestShrink_KeepKeys(t *testing.T) {
	es := baseShrinkSet(t)
	opts := DefaultShrinkOptions()
	opts.RemoveExpired = false
	opts.MaxKeys = 1
	opts.KeepKeys = []string{"ALPHA", "BETA"}

	if err := Shrink(es, opts); err == nil {
		t.Error("expected error when MaxKeys < len(KeepKeys)")
	}
}

func TestShrink_RemoveExpired(t *testing.T) {
	es := baseShrinkSet(t)
	past := time.Now().Add(-1 * time.Hour)
	if err := SetExpiry(es, "DELTA", past); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}

	opts := DefaultShrinkOptions()
	opts.RemoveExpired = true

	if err := Shrink(es, opts); err != nil {
		t.Fatalf("Shrink: %v", err)
	}
	if _, err := es.Get("DELTA"); err == nil {
		t.Error("expected DELTA to be removed (expired)")
	}
	if _, err := es.Get("ALPHA"); err != nil {
		t.Error("expected ALPHA to remain")
	}
}

func TestShrink_NilEnvSet(t *testing.T) {
	if err := Shrink(nil, DefaultShrinkOptions()); err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
