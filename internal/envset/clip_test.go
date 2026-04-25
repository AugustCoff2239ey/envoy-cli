package envset

import (
	"testing"
)

func baseClipSet() *EnvSet {
	es, _ := New("clipset", "test")
	es.Vars["ALPHA"] = "1"
	es.Vars["BETA"] = "2"
	es.Vars["GAMMA"] = "3"
	es.Vars["DELTA"] = "4"
	es.Vars["EPSILON"] = "5"
	return es
}

func TestClip_ReducesToMaxKeys(t *testing.T) {
	es := baseClipSet()
	opts := DefaultClipOptions()
	opts.MaxKeys = 3
	removed, err := Clip(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(es.Vars) != 3 {
		t.Errorf("expected 3 keys, got %d", len(es.Vars))
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(removed))
	}
}

func TestClip_ExplicitKeysRetained(t *testing.T) {
	es := baseClipSet()
	opts := DefaultClipOptions()
	opts.MaxKeys = 2
	opts.Keys = []string{"GAMMA", "EPSILON"}
	_, err := Clip(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := es.Vars["GAMMA"]; !ok {
		t.Error("expected GAMMA to be retained")
	}
	if _, ok := es.Vars["EPSILON"]; !ok {
		t.Error("expected EPSILON to be retained")
	}
}

func TestClip_LockedKeyNotRemoved(t *testing.T) {
	es := baseClipSet()
	_ = LockKey(es, "DELTA", "owner")
	opts := DefaultClipOptions()
	opts.MaxKeys = 2
	opts.SkipLocked = true
	_, err := Clip(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := es.Vars["DELTA"]; !ok {
		t.Error("expected locked key DELTA to be retained")
	}
}

func TestClip_NilEnvSet(t *testing.T) {
	_, err := Clip(nil, DefaultClipOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestClip_ZeroMaxKeys(t *testing.T) {
	es := baseClipSet()
	opts := DefaultClipOptions()
	opts.MaxKeys = 0
	_, err := Clip(es, opts)
	if err == nil {
		t.Error("expected error for MaxKeys=0")
	}
}
