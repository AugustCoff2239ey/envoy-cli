package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupRollbackStore(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	store = envset.NewStore(dir)
	rollbackStacks = map[string]*envset.RollbackStack{}

	es, _ := envset.New("myapp", "staging")
	es.Vars["API_KEY"] = "abc123"
	es.Vars["DB_URL"] = "postgres://localhost/dev"
	_ = store.Save(es)
}

func TestRollbackCmd_PushAndPop(t *testing.T) {
	setupRollbackStore(t)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"rollback", "push", "myapp", "staging", "-m", "before change"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("push failed: %v", err)
	}

	es, _ := store.Load("myapp", "staging")
	es.Vars["API_KEY"] = "changed"
	_ = store.Save(es)

	rootCmd.SetArgs([]string{"rollback", "pop", "myapp", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("pop failed: %v", err)
	}

	restored, err := store.Load("myapp", "staging")
	if err != nil {
		t.Fatalf("load after pop failed: %v", err)
	}
	if restored.Vars["API_KEY"] != "abc123" {
		t.Errorf("expected abc123, got %s", restored.Vars["API_KEY"])
	}
}

func TestRollbackCmd_NotFound(t *testing.T) {
	setupRollbackStore(t)

	rootCmd.SetArgs([]string{"rollback", "push", "ghost", "prod"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}

func TestRollbackCmd_PopEmpty(t *testing.T) {
	setupRollbackStore(t)

	rootCmd.SetArgs([]string{"rollback", "pop", "myapp", "staging"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error when popping empty stack")
	}
}
