package envset

import "testing"

func baseReadonlySet() *EnvSet {
	es, _ := New("readonly-test", "staging")
	_ = es.Set("API_KEY", "abc123")
	return es
}

func TestMarkReadonly_Valid(t *testing.T) {
	es := baseReadonlySet()
	if err := MarkReadonly(es); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsReadonly(es) {
		t.Error("expected envset to be read-only")
	}
}

func TestUnmarkReadonly_Valid(t *testing.T) {
	es := baseReadonlySet()
	_ = MarkReadonly(es)
	if err := UnmarkReadonly(es); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsReadonly(es) {
		t.Error("expected envset to not be read-only")
	}
}

func TestIsReadonly_Default(t *testing.T) {
	es := baseReadonlySet()
	if IsReadonly(es) {
		t.Error("new envset should not be read-only by default")
	}
}

func TestIsReadonly_NilEnvSet(t *testing.T) {
	if IsReadonly(nil) {
		t.Error("nil envset should not be read-only")
	}
}

func TestAssertWritable_Writable(t *testing.T) {
	es := baseReadonlySet()
	if err := AssertWritable(es); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestAssertWritable_Readonly(t *testing.T) {
	es := baseReadonlySet()
	_ = MarkReadonly(es)
	if err := AssertWritable(es); err == nil {
		t.Error("expected error for read-only envset")
	}
}

func TestMarkReadonly_NilEnvSet(t *testing.T) {
	if err := MarkReadonly(nil); err == nil {
		t.Error("expected error for nil envset")
	}
}
