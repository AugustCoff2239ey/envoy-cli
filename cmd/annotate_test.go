package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupAnnotateStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store := envset.NewStore(dir)
	e, _ := envset.New("myapp", "local")
	_ = e.Set("API_KEY", "secret")
	_ = e.Set("DB_URL", "postgres://localhost")
	_ = store.Save(e)
	return dir
}

func TestAnnotateCmd_AddAndList(t *testing.T) {
	dir := setupAnnotateStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"annotate", "API_KEY", "Primary key", "--name", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rootCmd.SetArgs([]string{"annotate", "API_KEY", "--name", "myapp", "--env", "local", "--list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
}

func TestAnnotateCmd_Remove(t *testing.T) {
	dir := setupAnnotateStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	rootCmd.SetArgs([]string{"annotate", "DB_URL", "some note", "--name", "myapp", "--env", "local"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"annotate", "DB_URL", "--name", "myapp", "--env", "local", "--remove"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnnotateCmd_NotFound(t *testing.T) {
	dir := setupAnnotateStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	rootCmd.SetArgs([]string{"annotate", "API_KEY", "note", "--name", "missing", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing EnvSet")
	}
}
