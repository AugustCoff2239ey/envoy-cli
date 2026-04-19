package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/envset"
)

func setupTypeCastStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/store.json"
	s, _ := envset.NewStore(path)
	es, _ := envset.New("myapp", "local")
	es.Vars["COUNT"] = "7.0"
	es.Vars["LABEL"] = "hello"
	es.Vars["FLAG"] = "1"
	_ = s.Save(es)
	return path
}

func TestTypeCastCmd_ToInt(t *testing.T) {
	path := setupTypeCastStore(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{
		"typecast", "--name", "myapp", "--env", "local",
		"--type", "int", "--key", "COUNT",
		"--store", path,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("COUNT")) {
		t.Errorf("expected COUNT in output, got: %s", buf.String())
	}
}

func TestTypeCastCmd_ToUpper(t *testing.T) {
	path := setupTypeCastStore(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{
		"typecast", "--name", "myapp", "--env", "local",
		"--type", "upper", "--key", "LABEL",
		"--store", path,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("HELLO")) {
		t.Errorf("expected HELLO in output, got: %s", buf.String())
	}
}

func TestTypeCastCmd_NotFound(t *testing.T) {
	path := setupTypeCastStore(t)
	rootCmd.SetArgs([]string{
		"typecast", "--name", "missing", "--env", "local",
		"--type", "int", "--store", path,
	})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing envset")
	}
}
