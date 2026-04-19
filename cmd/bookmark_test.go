package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupBookmarkStore(t *testing.T) *envset.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := envset.NewStore(dir)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	e, _ := envset.New("myapp", "local")
	_ = e.Set("DB_URL", "postgres://localhost/db")
	_ = st.Save(e)
	return st
}

func TestBookmarkCmd_AddAndList(t *testing.T) {
	st := setupBookmarkStore(t)
	overrideStore(t, st)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"bookmark", "add", "myapp", "prod-ref", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}

	rootCmd.SetArgs([]string{"bookmark", "list", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestBookmarkCmd_Remove(t *testing.T) {
	st := setupBookmarkStore(t)
	overrideStore(t, st)

	rootCmd.SetArgs([]string{"bookmark", "add", "myapp", "temp", "--env", "local"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"bookmark", "remove", "myapp", "temp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove: %v", err)
	}
}

func TestBookmarkCmd_NotFound(t *testing.T) {
	st := setupBookmarkStore(t)
	overrideStore(t, st)

	rootCmd.SetArgs([]string{"bookmark", "add", "ghost", "bm", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
