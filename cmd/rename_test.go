package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupRenameStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store, err := envset.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	es, err := envset.New("myapp", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("API_KEY", "secret")
	if err := store.Save(es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return dir
}

func TestRenameCmd_NewName(t *testing.T) {
	dir := setupRenameStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rename", "myapp", "staging", "--new-name", "webapp"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store, _ := envset.NewStore(dir)
	if _, err := store.Load("webapp", "staging"); err != nil {
		t.Errorf("expected renamed envset to exist: %v", err)
	}
	if _, err := store.Load("myapp", "staging"); err == nil {
		t.Error("expected old envset to be removed")
	}
}

func TestRenameCmd_NotFound(t *testing.T) {
	dir := setupRenameStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rename", "ghost", "production", "--new-name", "other"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for non-existent envset")
	}
}

func TestRenameCmd_NoFlags(t *testing.T) {
	dir := setupRenameStore(t)
	t.Setenv("ENVOY_STORE_DIR", dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rename", "myapp", "staging"})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error when no rename flags provided")
	}
}
