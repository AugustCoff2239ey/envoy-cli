package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupMirrorStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/store.json"

	store, _ := envset.NewStore(path)

	src, _ := envset.New("source", "local")
	src.Vars["API_KEY"] = "abc123"
	src.Vars["DB_HOST"] = "localhost"
	_ = store.Save(src)

	dst, _ := envset.New("dest", "local")
	dst.Vars["EXISTING"] = "keep"
	_ = store.Save(dst)

	return path
}

func TestMirrorCmd_AllKeys(t *testing.T) {
	path := setupMirrorStore(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{
		"mirror",
		"--store", path,
		"--src", "source",
		"--src-env", "local",
		"--dst", "dest",
		"--dst-env", "local",
	})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store, _ := envset.NewStore(path)
	dst, _ := store.Load("dest", "local")
	if dst.Vars["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY mirrored")
	}
	if dst.Vars["EXISTING"] != "keep" {
		t.Errorf("EXISTING should be preserved")
	}
}

func TestMirrorCmd_NotFound(t *testing.T) {
	path := setupMirrorStore(t)
	rootCmd.SetArgs([]string{
		"mirror",
		"--store", path,
		"--src", "ghost",
		"--src-env", "local",
		"--dst", "dest",
		"--dst-env", "local",
	})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestMirrorCmd_WithPrefix(t *testing.T) {
	path := setupMirrorStore(t)
	rootCmd.SetArgs([]string{
		"mirror",
		"--store", path,
		"--src", "source",
		"--src-env", "local",
		"--dst", "dest",
		"--dst-env", "local",
		"--prefix", "MIR_",
	})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store, _ := envset.NewStore(path)
	dst, _ := store.Load("dest", "local")
	if dst.Vars["MIR_API_KEY"] != "abc123" {
		t.Errorf("expected MIR_API_KEY to be set")
	}
}
