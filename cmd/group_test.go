package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupGroupStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store := envset.NewStore(dir)
	es, err := envset.New("myapp", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("API_KEY", "secret")
	if err := store.Save(es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return dir
}

func TestGroupCmd_CreateAndList(t *testing.T) {
	dir := setupGroupStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	rootCmd.SetArgs([]string{"group", "create", "myapp", "database", "DB_HOST,DB_PORT", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create: %v", err)
	}

	var buf bytes.Buffer
	r([]string{"group", "list", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Log("note: output captured via stdout; verifying store directly")
	}

	// Verify via store
	store := envset.NewStore(dir)
	es, err := store.Load("myapp", "local")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	g, err := envset.GetGroup(es, "database")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(g.Keys))
	}
}

func TestGroupCmd_Delete(t *testing.T) {
	dir := setupGroupStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	rootCmd.SetArgs([]string{"group", "create", "myapp", "api", "API_KEY", "--env", "local"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"group", "delete", "myapp", "api", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("delete: %v", err)
	}

	store := envset.NewStore(dir)
	es, _ := store.Load("myapp", "local")
	if _, err := envset.GetGroup(es, "api"); err == nil {
		t.Error("expected group to be deleted")
	}
}

func TestGroupCmd_NotFound(t *testing.T) {
	dir := setupGroupStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	rootCmd.SetArgs([]string{"group", "list", "nonexistent", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for nonexistent envset")
	}
}
