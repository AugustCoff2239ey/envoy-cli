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

func TestTagCmd_AddAndList(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"tag", "myapp", "add", "infra", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add tag: %v", err)
	}
	if !strings.Contains(buf.String(), "infra") {
		t.Errorf("expected confirmation mentioning 'infra', got: %s", buf.String())
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"tag", "myapp", "list", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list tags: %v", err)
	}
	if !strings.Contains(buf.String(), "infra") {
		t.Errorf("expected 'infra' in list output, got: %s", buf.String())
	}
}

func TestTagCmd_Remove(t *testing.T) {
	storePath := setupTagStore(t)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"tag", "myapp", "add", "to-remove", "--env", "local"})
	_ = rootCmd.Execute()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"tag", "myapp", "remove", "to-remove", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove tag: %v", err)
	}
	if !strings.Contains(buf.String(), "to-remove") {
		t.Errorf("expected confirmation mentioning 'to-remove', got: %s", buf.String())
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

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"tag", "myapp", "list", "--env", "local", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list json: %v", err)
	}
	if !strings.Contains(buf.String(), "json-tag") {
		t.Errorf("expected JSON output with 'json-tag', got: %s", buf.String())
	}
}
