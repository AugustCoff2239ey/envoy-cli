package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/envset"
)

func setupTemplateStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	f := filepath.Join(dir, "store.json")

	store, err := envset.NewStore(f)
	if err != nil {
		t.Fatal(err)
	}
	e, _ := envset.New("myapp", "local")
	e.Vars["DB_HOST"] = "localhost"
	e.Vars["DB_PORT"] = "5432"
	if err := store.Save(e); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestTemplateCmd_InlineResolved(t *testing.T) {
	f := setupTemplateStore(t)
	storeFile = f

	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{
		"template",
		"--name", "myapp",
		"--env", "local",
		"--template", "{{DB_HOST}}:{{DB_PORT}}",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplateCmd_FromFile(t *testing.T) {
	f := setupTemplateStore(t)
	storeFile = f

	tmplFile := filepath.Join(t.TempDir(), "tmpl.txt")
	_ = os.WriteFile(tmplFile, []byte("host={{DB_HOST}}"), 0o644)

	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{
		"template",
		"--name", "myapp",
		"--env", "local",
		"--file", tmplFile,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplateCmd_NotFound(t *testing.T) {
	f := setupTemplateStore(t)
	storeFile = f

	rootCmd.SetArgs([]string{
		"template",
		"--name", "ghost",
		"--env", "local",
		"--template", "{{DB_HOST}}",
	})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
