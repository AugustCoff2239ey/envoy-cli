package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envoy-cli/internal/envset"
)

func setupArchiveStore(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	store = envset.NewStore(dir)
	globalArchive = envset.NewArchive()
}

func TestArchiveCmd_AddAndList(t *testing.T) {
	setupArchiveStore(t)

	es, _ := envset.New("webapp", "staging")
	_ = es.Set("PORT", "8080")
	_ = store.Save(es)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"archive", "add", "webapp", "staging", "--reason", "test checkpoint"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("archive add: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"archive", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("archive list: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "test checkpoint") {
		t.Errorf("expected 'test checkpoint' in output, got: %s", out)
	}
}

func TestArchiveCmd_NotFound(t *testing.T) {
	setupArchiveStore(t)

	rootCmd.SetArgs([]string{"archive", "add", "missing", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}

func TestArchiveCmd_RestoreNotFound(t *testing.T) {
	setupArchiveStore(t)

	rootCmd.SetArgs([]string{"archive", "restore", "bad-id", "app", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for unknown archive ID")
	}
}
