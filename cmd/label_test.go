package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupLabelStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store, err := envset.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	es, _ := envset.New("myapp", "local")
	_ = store.Save(es)
	return dir
}

func TestLabelCmd_AddAndList(t *testing.T) {
	dir := setupLabelStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"label", "add", "--name", "myapp", "--env", "local", "team", "backend"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add label: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"label", "list", "--name", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list labels: %v", err)
	}
}

func TestLabelCmd_Remove(t *testing.T) {
	dir := setupLabelStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"label", "add", "--name", "myapp", "--env", "local", "owner", "alice"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"label", "remove", "--name", "myapp", "--env", "local", "owner"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove label: %v", err)
	}
}

func TestLabelCmd_NotFound(t *testing.T) {
	dir := setupLabelStore(t)
	t.Setenv("ENVOY_STORE", dir)

	rootCmd.SetArgs([]string{"label", "add", "--name", "ghost", "--env", "local", "k", "v"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
