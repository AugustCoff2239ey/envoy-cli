package envset

import (
	"testing"
)

func baseFmtSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("fmt-test", "development")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	es.Set("zebra_key", "  hello  ")
	es.Set("alpha_key", "  world  ")
	es.Set("middle_key", `say "hi"`)
	return es
}

func TestFmt_TrimValues(t *testing.T) {
	es := baseFmtSet(t)
	opts := DefaultFmtOptions()
	opts.SortKeys = false
	if err := Fmt(es, opts); err != nil {
		t.Fatalf("Fmt: %v", err)
	}
	v, _ := es.Get("zebra_key")
	if v != "hello" {
		t.Errorf("expected 'hello', got %q", v)
	}
}

func TestFmt_UppercaseKeys(t *testing.T) {
	es := baseFmtSet(t)
	opts := DefaultFmtOptions()
	opts.UppercaseKeys = true
	if err := Fmt(es, opts); err != nil {
		t.Fatalf("Fmt: %v", err)
	}
	if _, ok := es.Get("ALPHA_KEY"); !ok {
		t.Error("expected ALPHA_KEY to exist after uppercase transform")
	}
	if _, ok := es.Get("alpha_key"); ok {
		t.Error("original lowercase key should be gone")
	}
}

func TestFmt_SortKeys(t *testing.T) {
	es := baseFmtSet(t)
	opts := DefaultFmtOptions()
	if err := Fmt(es, opts); err != nil {
		t.Fatalf("Fmt: %v", err)
	}
	keys := es.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
			break
		}
	}
}

func TestFmt_QuoteValues(t *testing.T) {
	es := baseFmtSet(t)
	opts := FmtOptions{QuoteValues: true, TrimValues: false, SortKeys: false}
	if err := Fmt(es, opts); err != nil {
		t.Fatalf("Fmt: %v", err)
	}
	v, _ := es.Get("middle_key")
	expected := `"say \"hi\""`
	if v != expected {
		t.Errorf("expected %s, got %s", expected, v)
	}
}

func TestFmt_NilEnvSet(t *testing.T) {
	err := Fmt(nil, DefaultFmtOptions())
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestFmt_FrozenEnvSet(t *testing.T) {
	es := baseFmtSet(t)
	Freeze(es)
	err := Fmt(es, DefaultFmtOptions())
	if err == nil {
		t.Error("expected error for frozen EnvSet")
	}
}
