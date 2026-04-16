package cmd

import (
	"bytes"
	"testing"
	"time"

	"envoy-cli/internal/envset"
)

func setupExpireStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	store := envset.NewStore(dir + "/store.json")
	e, _ := envset.New("myapp", "local")
	e.Vars["TOKEN"] = "abc123"
	e.Vars["SECRET"] = "xyz"
	_ = store.Save(e)
	return dir + "/store.json"
}

func TestExpireCmd_SetAndGet(t *testing.T) {
	sf := setupExpireStore(t)
	t.Setenv("ENVOY_STORE", sf)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"expire", "TOKEN", "2h", "--name", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExpireCmd_Purge(t *testing.T) {
	sf := setupExpireStore(t)
	t.Setenv("ENVOY_STORE", sf)

	store := envset.NewStore(sf)
	e, _ := store.Load("myapp", "local")
	_ = envset.SetExpiry(e, "TOKEN", time.Now().Add(-time.Second))
	_ = store.Save(e)

	rootCmd.SetArgs([]string{"expire", "--name", "myapp", "--env", "local", "--purge"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e2, _ := store.Load("myapp", "local")
	if _, ok := e2.Vars["TOKEN"]; ok {
		t.Error("expected TOKEN to be purged")
	}
}

func TestExpireCmd_NotFound(t *testing.T) {
	sf := setupExpireStore(t)
	t.Setenv("ENVOY_STORE", sf)

	rootCmd.SetArgs([]string{"expire", "TOKEN", "1h", "--name", "ghost", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
