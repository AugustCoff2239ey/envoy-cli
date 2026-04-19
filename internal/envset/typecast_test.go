package envset

import (
	"testing"
)

func baseTypeCastSet() *EnvSet {
	es, _ := New("typecast-test", "local")
	es.Vars["COUNT"] = "42.0"
	es.Vars["RATIO"] = "3.14"
	es.Vars["ENABLED"] = "true"
	es.Vars["LABEL"] = "Hello"
	return es
}

func TestTypeCast_ToInt(t *testing.T) {
	es := baseTypeCastSet()
	results, err := TypeCast(es, "int", []string{"COUNT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Casted != "42" {
		t.Errorf("expected casted=42, got %v", results)
	}
	if es.Vars["COUNT"] != "42" {
		t.Errorf("expected envset updated, got %q", es.Vars["COUNT"])
	}
}

func TestTypeCast_ToBool(t *testing.T) {
	es := baseTypeCastSet()
	_, err := TypeCast(es, "bool", []string{"ENABLED"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["ENABLED"] != "true" {
		t.Errorf("expected true, got %q", es.Vars["ENABLED"])
	}
}

func TestTypeCast_ToUpper(t *testing.T) {
	es := baseTypeCastSet()
	_, err := TypeCast(es, "upper", []string{"LABEL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Vars["LABEL"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", es.Vars["LABEL"])
	}
}

func TestTypeCast_AllKeys(t *testing.T) {
	es, _ := New("tc", "local")
	es.Vars["A"] = "hello"
	es.Vars["B"] = "world"
	results, err := TypeCast(es, "upper", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestTypeCast_InvalidCast(t *testing.T) {
	es := baseTypeCastSet()
	_, err := TypeCast(es, "int", []string{"LABEL"})
	if err == nil {
		t.Error("expected error casting non-numeric to int")
	}
}

func TestTypeCast_UnsupportedType(t *testing.T) {
	es := baseTypeCastSet()
	_, err := TypeCast(es, "base64", []string{"LABEL"})
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestTypeCast_NilEnvSet(t *testing.T) {
	_, err := TypeCast(nil, "int", nil)
	if err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestTypeCast_MissingKey(t *testing.T) {
	es := baseTypeCastSet()
	_, err := TypeCast(es, "int", []string{"NONEXISTENT"})
	if err == nil {
		t.Error("expected error for missing key")
	}
}
