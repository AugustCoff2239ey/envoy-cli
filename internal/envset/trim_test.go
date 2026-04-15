package envset

import (
	"testing"
)

func baseTrimSet() *EnvSet {
	es, _ := New("trim-test", "local")
	es.Vars["CLEAN"] = "already_clean"
	es.Vars["PADDED"] = "  hello world  "
	es.Vars["TAB_PAD"] = "\tvalue\t"
	es.Vars["QUOTED_D"] = `"quoted"`
	es.Vars["QUOTED_S"] = `'single'`
	return es
}

func TestTrim_WhitespaceDefault(t *testing.T) {
	es := baseTrimSet()
	opts := DefaultTrimOptions()
	results, err := Trim(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changed := map[string]TrimResult{}
	for _, r := range results {
		if r.Changed {
			changed[r.Key] = r
		}
	}
	if _, ok := changed["PADDED"]; !ok {
		t.Error("expected PADDED to be changed")
	}
	if es.Vars["PADDED"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", es.Vars["PADDED"])
	}
	if _, ok := changed["CLEAN"]; ok {
		t.Error("CLEAN should not be changed")
	}
}

func TestTrim_Quotes(t *testing.T) {
	es := baseTrimSet()
	opts := DefaultTrimOptions()
	opts.Quotes = true
	_, err := Trim(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["QUOTED_D"] != "quoted" {
		t.Errorf("expected 'quoted', got %q", es.Vars["QUOTED_D"])
	}
	if es.Vars["QUOTED_S"] != "single" {
		t.Errorf("expected 'single', got %q", es.Vars["QUOTED_S"])
	}
}

func TestTrim_SelectedKeys(t *testing.T) {
	es := baseTrimSet()
	opts := DefaultTrimOptions()
	opts.Keys = []string{"PADDED"}
	_, err := Trim(es, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["PADDED"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", es.Vars["PADDED"])
	}
	if es.Vars["TAB_PAD"] != "\tvalue\t" {
		t.Error("TAB_PAD should be untouched when not in Keys list")
	}
}

func TestTrim_MissingKey(t *testing.T) {
	es := baseTrimSet()
	opts := DefaultTrimOptions()
	opts.Keys = []string{"DOES_NOT_EXIST"}
	_, err := Trim(es, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestTrim_NilEnvSet(t *testing.T) {
	_, err := Trim(nil, DefaultTrimOptions())
	if err == nil {
		t.Fatal("expected error for nil envset")
	}
}
