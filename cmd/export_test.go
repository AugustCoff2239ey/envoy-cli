package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/envset"
)

func setupExportStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store, err := envset.NewStore(filepath.Join(dir, "envsets"))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	set, err := envset.New("myapp", "staging")
	if err != nil {
		t.Fatalf("failed to create env set: %v", err)
	}
	set.Set("API_KEY", "abc123")
	set.Set("DEBUG", "true")

	if err := store.Save(set); err != nil {
		t.Fatalf("failed to save env set: %v", err)
	}
	return dir
}

func TestExportCmd_Dotenv(t *testing.T) {
	dir := setupExportStore(t)
	t.Setenv("ENVOY_STORE_PATH", filepath.Join(dir, "envsets"))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"export", "myapp", "staging", "--format", "dotenv"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "API_KEY=abc123") {
		t.Errorf("expected dotenv output to contain API_KEY=abc123, got: %s", out)
	}
}

func TestExportCmd_ToFile(t *testing.T) {
	dir := setupExportStore(t)
	t.Setenv("ENVOY_STORE_PATH", filepath.Join(dir, "envsets"))

	outFile := filepath.Join(dir, ".env")
	rootCmd.SetArgs([]string{"export", "myapp", "staging", "--format", "dotenv", "--output", outFile})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
	if !strings.Contains(string(data), "API_KEY=abc123") {
		t.Errorf("expected file to contain API_KEY=abc123, got: %s", string(data))
	}
}

func TestExportCmd_NotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVOY_STORE_PATH", filepath.Join(dir, "envsets"))

	rootCmd.SetArgs([]string{"export", "ghost", "production"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing env set, got nil")
	}
}
