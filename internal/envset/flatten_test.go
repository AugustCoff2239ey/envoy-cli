package envset

import (
	"testing"
)

func baseFlattenSets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	a, _ := New("base", "local")
	a.Vars["HOST"] = "localhost"
	a.Vars["PORT"] = "8080"
	a.Vars["DEBUG"] = "true"

	b, _ := New("override", "local")
	b.Vars["PORT"] = "9090"
	b.Vars["NEW_KEY"] = "hello"

	return a, b
}

func TestFlatten_NoConflict(t *testing.T) {
	a, _ := New("a", "local")
	a.Vars["FOO"] = "1"

	b, _ := New("b", "local")
	b.Vars["BAR"] = "2"

	res, err := Flatten(FlattenOptions{}, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "1" || res.Vars["BAR"] != "2" {
		t.Errorf("expected merged vars, got %v", res.Vars)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestFlatten_ConflictNoOverwrite(t *testing.T) {
	a, b := baseFlattenSets(t)
	res, err := Flatten(FlattenOptions{Overwrite: false}, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["PORT"] != "8080" {
		t.Errorf("expected original PORT=8080, got %s", res.Vars["PORT"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "PORT" {
		t.Errorf("expected conflict on PORT, got %v", res.Conflicts)
	}
}

func TestFlatten_ConflictWithOverwrite(t *testing.T) {
	a, b := baseFlattenSets(t)
	res, err := Flatten(FlattenOptions{Overwrite: true}, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["PORT"] != "9090" {
		t.Errorf("expected overwritten PORT=9090, got %s", res.Vars["PORT"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	a, _ := New("a", "local")
	a.Vars["KEY"] = "val"

	res, err := Flatten(FlattenOptions{Prefix: "APP_"}, a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_KEY"] != "val" {
		t.Errorf("expected prefixed key APP_KEY, got %v", res.Vars)
	}
}

func TestFlatten_UppercaseKeys(t *testing.T) {
	a, _ := New("a", "local")
	a.Vars["my_key"] = "value"

	res, err := Flatten(FlattenOptions{UppercaseKeys: true}, a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["MY_KEY"] != "value" {
		t.Errorf("expected MY_KEY, got %v", res.Vars)
	}
}

func TestFlatten_NilSource(t *testing.T) {
	a, _ := New("a", "local")
	_, err := Flatten(FlattenOptions{}, a, nil)
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestFlatten_EmptySets(t *testing.T) {
	res, err := Flatten(FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty vars, got %v", res.Vars)
	}
}
