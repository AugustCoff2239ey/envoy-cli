package envset

import (
	"testing"
)

func baseRollbackSet() *EnvSet {
	es, _ := New("rollback-test", "staging")
	es.Vars["DB_HOST"] = "localhost"
	es.Vars["DB_PORT"] = "5432"
	return es
}

func TestRollback_PushAndPop(t *testing.T) {
	es := baseRollbackSet()
	stack := NewRollbackStack(5)

	if err := stack.Push(es, "initial state"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	es.Vars["DB_HOST"] = "remotehost"

	if err := stack.Pop(es); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if es.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", es.Vars["DB_HOST"])
	}
}

func TestRollback_PopEmptyStack(t *testing.T) {
	es := baseRollbackSet()
	stack := NewRollbackStack(5)

	if err := stack.Pop(es); err == nil {
		t.Error("expected error on empty stack pop")
	}
}

func TestRollback_MaxCapacity(t *testing.T) {
	es := baseRollbackSet()
	stack := NewRollbackStack(3)

	for i := 0; i < 5; i++ {
		_ = stack.Push(es, "state")
	}

	if stack.Len() != 3 {
		t.Errorf("expected max 3 entries, got %d", stack.Len())
	}
}

func TestRollback_NilEnvSet(t *testing.T) {
	stack := NewRollbackStack(5)

	if err := stack.Push(nil, "msg"); err == nil {
		t.Error("expected error for nil envset on push")
	}
	if err := stack.Pop(nil); err == nil {
		t.Error("expected error for nil envset on pop")
	}
}

func TestRollback_Peek(t *testing.T) {
	es := baseRollbackSet()
	stack := NewRollbackStack(5)

	_ = stack.Push(es, "before change")

	entry, err := stack.Peek()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Message != "before change" {
		t.Errorf("expected 'before change', got %s", entry.Message)
	}
	if stack.Len() != 1 {
		t.Errorf("peek should not remove entry, len=%d", stack.Len())
	}
}

func TestRollback_Independence(t *testing.T) {
	es := baseRollbackSet()
	stack := NewRollbackStack(5)
	_ = stack.Push(es, "snapshot")

	es.Vars["NEW_KEY"] = "new_value"

	entry, _ := stack.Peek()
	if _, ok := entry.Vars["NEW_KEY"]; ok {
		t.Error("snapshot should be independent of later mutations")
	}
}
