package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envset"
)

func setupTagStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store, err := envset.NewStore(filepath.Join(dir, "store.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	es, _ := envset.New("myapp", "local")
	if err := store.Save(es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return filepath.Join(dir, "store.json")
}

// runTagCmd executes the root command with the given args and returns the output.
func runTagCmd(t *testing.T, args []string) string {
	t.Helper()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("command %v: %v", args, err)
	}
	return buf.String()
}

func TestTagCmd_AddAndList(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	out := runTagCmd(t, []string{"tag", "myapp", "add", "infra", "--env", "local"})
	if !strings.Contains(out, "infra") {
		t.Errorf("expected confirmation mentioning 'infra', got: %s", out)
	}

	out = runTagCmd(t, []string{"tag", "myapp", "list", "--env", "local"})
	if !strings.Contains(out, "infra") {
		t.Errorf("expected 'infra' in list output, got: %s", out)
	}
}

func TestTagCmd_Remove(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"tag", "myapp", "add", "to-remove", "--env", "local"})
	_ = rootCmd.Execute()

	out := runTagCmd(t, []string{"tag", "myapp", "remove", "to-remove", "--env", "local"})
	if !strings.Contains(out, "to-remove") {
		t.Errorf("expected confirmation mentioning 'to-remove', got: %s", out)
	}
}

func TestTagCmd_NotFound(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"tag", "ghost", "list", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for non-existent envset")
	}
}

func TestTagCmd_ListJSON(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"tag", "myapp", "add", "json-tag", "--env", "local"})
	_ = rootCmd.Execute()

	out := runTagCmd(t, []string{"tag", "myapp", "list", "--env", "local", "--json"})
	if !strings.Contains(out, "json-tag") {
		t.Errorf("expected JSON output with 'json-tag', got: %s", out)
	}
}
