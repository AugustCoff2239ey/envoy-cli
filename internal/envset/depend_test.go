package envset

import (
	"testing"
)

func baseDependSet() *EnvSet {
	es, _ := New("deps", "local")
	es.Vars["DB_HOST"] = "localhost"
	es.Vars["DB_PORT"] = "5432"
	es.Vars["DB_URL"] = "postgres://localhost:5432"
	return es
}

func TestAddDependency_Valid(t *testing.T) {
	es := baseDependSet()
	if err := AddDependency(es, "DB_URL", "DB_HOST"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := AddDependency(es, "DB_URL", "DB_PORT"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps := GetDependencies(es, "DB_URL")
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d", len(deps))
	}
}

func TestAddDependency_Duplicate(t *testing.T) {
	es := baseDependSet()
	_ = AddDependency(es, "DB_URL", "DB_HOST")
	_ = AddDependency(es, "DB_URL", "DB_HOST")
	if len(GetDependencies(es, "DB_URL")) != 1 {
		t.Fatal("expected dedup of duplicate dependency")
	}
}

func TestAddDependency_MissingKey(t *testing.T) {
	es := baseDependSet()
	if err := AddDependency(es, "MISSING", "DB_HOST"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAddDependency_NilEnvSet(t *testing.T) {
	if err := AddDependency(nil, "DB_URL", "DB_HOST"); err == nil {
		t.Fatal("expected error for nil EnvSet")
	}
}

func TestRemoveDependency_Valid(t *testing.T) {
	es := baseDependSet()
	_ = AddDependency(es, "DB_URL", "DB_HOST")
	_ = AddDependency(es, "DB_URL", "DB_PORT")
	_ = RemoveDependency(es, "DB_URL", "DB_HOST")
	deps := GetDependencies(es, "DB_URL")
	if len(deps) != 1 || deps[0] != "DB_PORT" {
		t.Fatalf("unexpected deps after remove: %v", deps)
	}
}

func TestCheckDependencies_AllPresent(t *testing.T) {
	es := baseDependSet()
	_ = AddDependency(es, "DB_URL", "DB_HOST")
	missing := CheckDependencies(es)
	if len(missing) != 0 {
		t.Fatalf("expected no missing deps, got %v", missing)
	}
}

func TestCheckDependencies_MissingDep(t *testing.T) {
	es := baseDependSet()
	_ = AddDependency(es, "DB_URL", "DB_HOST")
	delete(es.Vars, "DB_HOST")
	missing := CheckDependencies(es)
	if len(missing) != 1 || missing[0] != "DB_HOST" {
		t.Fatalf("expected DB_HOST missing, got %v", missing)
	}
}
