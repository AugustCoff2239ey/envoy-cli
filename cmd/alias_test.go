package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupAliasStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store, err := envset.NewStore(dir)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	es, _ := envset.New("myapp", "local")
	_ = es.Set("DB_HOST", "localhost")
	if err := store.Save(es); err != nil {
		t.Fatalf("save: %v", err)
	}
	return dir
}

func TestAliasCmd_AddAndList(t *testing.T) {
	dir := setupAliasStore(t)
	t.Setenv("ENVOY_STORE", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"alias", "add", "DB_HOST", "db-host", "--name", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add alias: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"alias", "list", "DB_HOST", "--name", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list alias: %v", err)
	}
	if out := buf.String(); out == "" {
		t.Error("expected alias output")
	}
}

func TestAliasCmd_Remove(t *testing.T) {
	dir := setupAliasStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"alias", "add", "DB_HOST", "host", "--name", "myapp"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"alias", "remove", "DB_HOST", "host", "--name", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove alias: %v", err)
	}
}

func TestAliasCmd_NotFound(t *testing.T) {
	dir := setupAliasStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"alias", "add", "DB_HOST", "host", "--name", "ghost"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
