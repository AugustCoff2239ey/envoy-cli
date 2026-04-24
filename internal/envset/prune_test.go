package envset

import (
	"testing"
)

func basePruneSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("prune-test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Vars["HOST"] = "localhost"
	es.Vars["PORT"] = ""
	es.Vars["DEBUG"] = "true"
	es.Vars["VERBOSE"] = "true" // duplicate value of DEBUG
	es.Vars["EMPTY"] = ""
	return es
}

func TestPrune_RemoveEmpty(t *testing.T) {
	es := basePruneSet(t)
	opts := DefaultPruneOptions()
	n, err := Prune(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 removed, got %d", n)
	}
	if _, ok := es.Vars["PORT"]; ok {
		t.Error("PORT should have been pruned")
	}
	if _, ok := es.Vars["EMPTY"]; ok {
		t.Error("EMPTY should have been pruned")
	}
}

func TestPrune_RemoveDuplicateValues(t *testing.T) {
	es := basePruneSet(t)
	opts := PruneOptions{RemoveEmpty: false, RemoveDuplicateValues: true}
	n, err := Prune(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DEBUG and VERBOSE share value "true"; one should be removed.
	if n != 1 {
		t.Errorf("expected 1 removed, got %d", n)
	}
}

func TestPrune_ExplicitKeys(t *testing.T) {
	es := basePruneSet(t)
	opts := PruneOptions{Keys: []string{"HOST", "DEBUG"}}
	n, err := Prune(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 removed, got %d", n)
	}
	if _, ok := es.Vars["HOST"]; ok {
		t.Error("HOST should have been pruned")
	}
}

func TestPrune_LockedKeyBlocked(t *testing.T) {
	es := basePruneSet(t)
	if err := LockKey(es, "PORT", "test"); err != nil {
		t.Fatalf("LockKey: %v", err)
	}
	_, err := Prune(es, DefaultPruneOptions())
	if err == nil {
		t.Error("expected error when pruning locked key")
	}
}

func TestPrune_NilEnvSet(t *testing.T) {
	_, err := Prune(nil, DefaultPruneOptions())
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestPrune_ExplicitKeyMissing(t *testing.T) {
	es := basePruneSet(t)
	opts := PruneOptions{Keys: []string{"NONEXISTENT"}}
	n, err := Prune(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 removed, got %d", n)
	}
}
