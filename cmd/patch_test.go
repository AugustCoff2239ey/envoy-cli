package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupPatchStore(t *testing.T) *envset.Store {
	t.Helper()
	dir := t.TempDir()
	s := envset.NewStore(dir)
	es, _ := envset.New("myapp", "staging")
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("DEBUG", "true")
	_ = s.Save(es)
	return s
}

func TestPatchCmd_SetAndDelete(t *testing.T) {
	s := setupPatchStore(t)
	overrideStore(t, s)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{
		"patch", "myapp", "staging",
		"--set", "DB_PORT=5433",
		"--delete", "DEBUG",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	es, err := s.Load("myapp", "staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if v, _ := es.Get("DB_PORT"); v != "5433" {
		t.Errorf("expected DB_PORT=5433, got %q", v)
	}
	if _, ok := es.Get("DEBUG"); ok {
		t.Error("expected DEBUG to be deleted")
	}
}

func TestPatchCmd_NotFound(t *testing.T) {
	s := setupPatchStore(t)
	overrideStore(t, s)

	rootCmd.SetArgs([]string{"patch", "ghost", "prod", "--set", "X=1"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing env set")
	}
}

func TestPatchCmd_NoOps(t *testing.T) {
	s := setupPatchStore(t)
	overrideStore(t, s)

	rootCmd.SetArgs([]string{"patch", "myapp", "staging"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error when no operations provided")
	}
}
