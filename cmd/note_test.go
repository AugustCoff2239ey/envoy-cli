package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupNoteStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/store.json"
	st, err := envset.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	es, _ := envset.New("myapp", "staging")
	_ = es.Set("KEY", "val")
	if err := st.Save(es); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNoteCmd_AddAndList(t *testing.T) {
	path := setupNoteStore(t)
	storeFile = path

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"note", "add", "Deployed to staging", "--name", "myapp", "--env", "staging", "--author", "alice"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add note failed: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"note", "list", "--name", "myapp", "--env", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list notes failed: %v", err)
	}
}

func TestNoteCmd_Clear(t *testing.T) {
	path := setupNoteStore(t)
	storeFile = path

	rootCmd.SetArgs([]string{"note", "add", "Temp note", "--name", "myapp", "--env", "staging"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"note", "clear", "--name", "myapp", "--env", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("clear notes failed: %v", err)
	}

	st, _ := envset.NewStore(path)
	es, _ := st.Load("myapp", "staging")
	if notes := envset.ListNotes(es); len(notes) != 0 {
		t.Errorf("expected 0 notes after clear, got %d", len(notes))
	}
}

func TestNoteCmd_NotFound(t *testing.T) {
	path := setupNoteStore(t)
	storeFile = path

	rootCmd.SetArgs([]string{"note", "list", "--name", "ghost", "--env", "staging"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing env set")
	}
}
