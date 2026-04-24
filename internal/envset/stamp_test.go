package envset

import (
	"strings"
	"testing"
	"time"
)

func baseStampSet() *EnvSet {
	es, _ := New("stamp-test", "local")
	_ = es.Set("BUILT_AT", "old")
	_ = es.Set("DEPLOYED_AT", "old")
	_ = es.Set("NOTES", "unchanged")
	return es
}

func TestStamp_AllKeys(t *testing.T) {
	es := baseStampSet()
	opts := DefaultStampOptions()
	stamped, err := Stamp(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stamped) != 3 {
		t.Errorf("expected 3 stamped keys, got %d", len(stamped))
	}
	// Verify values parse as RFC3339.
	for k, v := range stamped {
		if _, err := time.Parse(time.RFC3339, v); err != nil {
			t.Errorf("key %q value %q is not valid RFC3339: %v", k, v, err)
		}
	}
}

func TestStamp_SelectedKeys(t *testing.T) {
	es := baseStampSet()
	opts := DefaultStampOptions()
	opts.Keys = []string{"BUILT_AT"}
	stamped, err := Stamp(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stamped) != 1 {
		t.Errorf("expected 1 stamped key, got %d", len(stamped))
	}
	if es.Vars["DEPLOYED_AT"] != "old" {
		t.Errorf("DEPLOYED_AT should remain unchanged")
	}
}

func TestStamp_PrefixSuffix(t *testing.T) {
	es := baseStampSet()
	opts := DefaultStampOptions()
	opts.Keys = []string{"BUILT_AT"}
	opts.Prefix = "ts:"
	opts.Suffix = ":end"
	_, err := Stamp(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v := es.Vars["BUILT_AT"]
	if !strings.HasPrefix(v, "ts:") {
		t.Errorf("expected prefix 'ts:', got %q", v)
	}
	if !strings.HasSuffix(v, ":end") {
		t.Errorf("expected suffix ':end', got %q", v)
	}
}

func TestStamp_SkipsLockedKey(t *testing.T) {
	es := baseStampSet()
	_ = LockKey(es, "BUILT_AT", "test")
	opts := DefaultStampOptions()
	stamped, err := Stamp(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := stamped["BUILT_AT"]; ok {
		t.Errorf("locked key BUILT_AT should have been skipped")
	}
}

func TestStamp_MissingKey(t *testing.T) {
	es := baseStampSet()
	opts := DefaultStampOptions()
	opts.Keys = []string{"NONEXISTENT"}
	_, err := Stamp(es, opts)
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestStamp_NilEnvSet(t *testing.T) {
	_, err := Stamp(nil, DefaultStampOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}
