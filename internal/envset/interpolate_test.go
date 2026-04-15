package envset

import (
	"testing"
)

func baseInterpolateSet() *EnvSet {
	e, _ := New("interpolate-test", "local")
	e.Vars["HOST"] = "localhost"
	e.Vars["PORT"] = "5432"
	e.Vars["DB_URL"] = "postgres://${HOST}:${PORT}/mydb"
	e.Vars["GREETING"] = "hello ${NAME:world}"
	e.Vars["MISSING_REF"] = "value is ${UNKNOWN}"
	return e
}

func TestInterpolate_ResolvesKnownRefs(t *testing.T) {
	e := baseInterpolateSet()
	res, err := Interpolate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Resolved["DB_URL"]
	want := "postgres://localhost:5432/mydb"
	if got != want {
		t.Errorf("DB_URL: got %q, want %q", got, want)
	}
}

func TestInterpolate_UsesDefault(t *testing.T) {
	e := baseInterpolateSet()
	res, err := Interpolate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Resolved["GREETING"]
	want := "hello world"
	if got != want {
		t.Errorf("GREETING: got %q, want %q", got, want)
	}
}

func TestInterpolate_TracksUnresolved(t *testing.T) {
	e := baseInterpolateSet()
	res, err := Interpolate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) == 0 {
		t.Fatal("expected at least one unresolved key")
	}
	found := false
	for _, u := range res.Unresolved {
		if u == "UNKNOWN" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected UNKNOWN in unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolate_PlainValueUnchanged(t *testing.T) {
	e := baseInterpolateSet()
	res, err := Interpolate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resolved["HOST"] != "localhost" {
		t.Errorf("HOST should remain localhost, got %q", res.Resolved["HOST"])
	}
}

func TestInterpolate_NilEnvSet(t *testing.T) {
	_, err := Interpolate(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestInterpolate_DeduplicatesUnresolved(t *testing.T) {
	e, _ := New("dup-test", "local")
	e.Vars["A"] = "${MISSING}"
	e.Vars["B"] = "${MISSING}"
	res, err := Interpolate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := 0
	for _, u := range res.Unresolved {
		if u == "MISSING" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected MISSING once in unresolved, got %d times", count)
	}
}
