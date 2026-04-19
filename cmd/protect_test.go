package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupProtectStore(t *testing.T) *envset.Store {
	t.Helper()
	s := tempTestStore(t)
	e, _ := envset.New("myapp", "local")
	_ = e.Set("API_KEY", "abc123")
	_ = e.Set("DB_PASS", "secret")
	_ = s.Save(e)
	return s
}

func TestProtectCmd_ProtectAndList(t *testing.T) {
	s := setupProtectStore(t)
	overridStore(t, s)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"protect", "myapp", "API_KEY", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("protect failed: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"list-protected", "myapp", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list-protected failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("API_KEY")) {
		t.Errorf("expected API_KEY in output, got: %s", buf.String())
	}
}

func TestProtectCmd_Unprotect(t *testing.T) {
	s := setupProtectStore(t)
	overridStore(t, s)

	rootCmd.SetArgs([]string{"protect", "myapp", "DB_PASS", "--env", "local"})
	_ = rootCmd.Execute()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"unprotect", "myapp", "DB_PASS", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unprotect failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("unprotected")) {
		t.Errorf("expected unprotected in output, got: %s", buf.String())
	}
}

func TestProtectCmd_NotFound(t *testing.T) {
	s := setupProtectStore(t)
	overridStore(t, s)

	rootCmd.SetArgs([]string{"protect", "ghost", "API_KEY", "--env", "local"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}

func TestProtectCmd_AlreadyProtected(t *testing.T) {
	s := setupProtectStore(t)
	overridStore(t, s)

	// Protect the key once
	rootCmd.SetArgs([]string{"protect", "myapp", "API_KEY", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first protect failed: %v", err)
	}

	// Protecting the same key again should not return an error
	rootCmd.SetArgs([]string{"protect", "myapp", "API_KEY", "--env", "local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("second protect failed: %v", err)
	}
}
