package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupScopeStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store := envset.NewStore(dir)
	es, err := envset.New("myapp", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, kv := range []struct{ k, v string }{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"API_KEY", "topsecret"},
	} {
		if err := es.Set(kv.k, kv.v); err != nil {
			t.Fatalf("Set %s: %v", kv.k, err)
		}
	}
	if err := store.Save(es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return dir
}

func TestScopeCmd_CreateAndList(t *testing.T) {
	dir := setupScopeStore(t)
	t.Setenv("ENVOY_STORE", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"scope", "create", "database", "myapp", "DB_HOST,DB_PORT"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create scope: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"scope", "list", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list scopes: %v", err)
	}
	if out := buf.String(); out == "" {
		t.Error("expected non-empty scope list output")
	}
}

func TestScopeCmd_Delete(t *testing.T) {
	dir := setupScopeStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"scope", "create", "api", "myapp", "API_KEY"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create scope: %v", err)
	}

	rootCmd.SetArgs([]string{"scope", "delete", "api", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("delete scope: %v", err)
	}
}

func TestScopeCmd_NotFound(t *testing.T) {
	dir := setupScopeStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"scope", "list", "ghost-set"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for non-existent env set")
	}
}

func TestScopeCmd_ListJSON(t *testing.T) {
	dir := setupScopeStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"scope", "create", "db", "myapp", "DB_HOST"})
	_ = rootCmd.Execute()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"scope", "list", "myapp", "--format", "json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list json: %v", err)
	}
	if out := buf.String(); out == "" {
		t.Error("expected JSON output")
	}
}
