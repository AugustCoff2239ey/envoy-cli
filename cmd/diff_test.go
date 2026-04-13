package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envset"
)

func setupDiffStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.db")

	store, err := envset.NewStore(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	base, _ := envset.New("local", "local")
	_ = base.Set("KEY_A", "value_a")
	_ = base.Set("KEY_B", "old_value")
	_ = store.Save(base)

	target, _ := envset.New("staging", "staging")
	_ = target.Set("KEY_A", "value_a")
	_ = target.Set("KEY_B", "new_value")
	_ = target.Set("KEY_C", "added")
	_ = store.Save(target)

	return path
}

func TestDiffCmd_TextOutput(t *testing.T) {
	path := setupDiffStore(t)
	old := os.Getenv("ENVOY_STORE")
	os.Setenv("ENVOY_STORE", path)
	defer os.Setenv("ENVOY_STORE", old)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"diff", "local", "staging"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "local → staging") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "+ KEY_C") {
		t.Errorf("expected added key in output, got: %s", out)
	}
	if !strings.Contains(out, "~ KEY_B") {
		t.Errorf("expected changed key in output, got: %s", out)
	}
}

func TestDiffCmd_NoDiff(t *testing.T) {
	path := setupDiffStore(t)
	old := os.Getenv("ENVOY_STORE")
	os.Setenv("ENVOY_STORE", path)
	defer os.Setenv("ENVOY_STORE", old)

	store, _ := envset.NewStore(path)
	same, _ := envset.New("staging2", "staging")
	_ = same.Set("KEY_A", "value_a")
	_ = same.Set("KEY_B", "old_value")
	_ = store.Save(same)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"diff", "local", "staging2"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}
