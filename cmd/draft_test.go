package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupDraftStore(t *testing.T) string {
	t.Helper()
	path := t.TempDir() + "/draft_test.json"
	store, err := envset.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	e, _ := envset.New("webapp", "staging")
	_ = e.Set("API_KEY", "abc123")
	_ = store.Save(e)
	return path
}

func TestDraftCmd_SaveAndPromote(t *testing.T) {
	path := setupDraftStore(t)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--store", path, "draft", "save", "webapp", "staging", "--note", "wip"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("save draft failed: %v", err)
	}

	store, _ := envset.NewStore(path)
	d, err := store.Load("draft:webapp", "staging")
	if err != nil {
		t.Fatalf("draft not found after save: %v", err)
	}
	if d.Meta["draft_note"] != "wip" {
		t.Errorf("expected note 'wip', got %q", d.Meta["draft_note"])
	}

	rootCmd.SetArgs([]string{"--store", path, "draft", "promote", "webapp", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("promote draft failed: %v", err)
	}

	promoted, err := store.Load("webapp", "staging")
	if err != nil {
		t.Fatalf("promoted envset not found: %v", err)
	}
	if _, ok := promoted.Meta["draft_note"]; ok {
		t.Error("draft_note should be absent after promotion")
	}
}

func TestDraftCmd_NotFound(t *testing.T) {
	path := setupDraftStore(t)
	rootCmd.SetArgs([]string{"--store", path, "draft", "save", "missing", "production"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
