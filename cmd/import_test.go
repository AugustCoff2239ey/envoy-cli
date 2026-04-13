package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envset"
)

func setupImportStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVOY_STORE_PATH", filepath.Join(dir, "store.json"))
	return dir
}

func TestImportCmd_FromFile(t *testing.T) {
	dir := setupImportStore(t)

	envFile := filepath.Join(dir, "test.env")
	content := "APP_KEY=abc123\nDEBUG=true\n"
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"import", "myapp", "--env", "local", "--file", envFile, "--format", "dotenv"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected output to mention envset name, got: %s", out)
	}

	store, _ := envset.NewStore("")
	set, err := store.Load("myapp", "local")
	if err != nil {
		t.Fatalf("failed to load saved envset: %v", err)
	}
	val, ok := set.Get("APP_KEY")
	if !ok || val != "abc123" {
		t.Errorf("expected APP_KEY=abc123, got %q (found=%v)", val, ok)
	}
}

func TestImportCmd_NotFound(t *testing.T) {
	setupImportStore(t)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"import", "myapp", "--env", "local", "--file", "/nonexistent/path.env"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestImportCmd_OverwriteFlag(t *testing.T) {
	dir := setupImportStore(t)

	store, _ := envset.NewStore("")
	initial, _ := envset.New("myapp", "local")
	_ = initial.Set("OLD_KEY", "old")
	_ = store.Save(initial)

	envFile := filepath.Join(dir, "new.env")
	_ = os.WriteFile(envFile, []byte("NEW_KEY=newval\n"), 0644)

	rootCmd.SetArgs([]string{"import", "myapp", "--env", "local", "--file", envFile, "--overwrite"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	set, err := store.Load("myapp", "local")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	if _, ok := set.Get("OLD_KEY"); ok {
		t.Error("expected OLD_KEY to be absent after overwrite")
	}
	if val, ok := set.Get("NEW_KEY"); !ok || val != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %q (found=%v)", val, ok)
	}
}
