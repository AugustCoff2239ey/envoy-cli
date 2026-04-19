package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupCompareStore(t *testing.T) string {
	t.Helper()
	f := tempStoreFile(t)
	store, _ := envset.NewStore(f)

	base, _ := envset.New("app", "staging")
	_ = base.Set("HOST", "localhost")
	_ = base.Set("PORT", "8080")
	_ = base.Set("ONLY_BASE", "yes")
	_ = store.Save(base)

	target, _ := envset.New("app", "production")
	_ = target.Set("HOST", "prod.example.com")
	_ = target.Set("PORT", "8080")
	_ = target.Set("ONLY_TARGET", "yes")
	_ = store.Save(target)

	return f
}

func TestCompareCmd_TextOutput(t *testing.T) {
	f := setupCompareStore(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--store", f, "compare", "app", "staging", "app", "production"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompareCmd_JSONOutput(t *testing.T) {
	f := setupCompareStore(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--store", f, "compare", "app", "staging", "app", "production", "-o", "json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompareCmd_NotFound(t *testing.T) {
	f := setupCompareStore(t)
	rootCmd.SetArgs([]string{"--store", f, "compare", "missing", "env", "app", "production"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing base set")
	}
}
