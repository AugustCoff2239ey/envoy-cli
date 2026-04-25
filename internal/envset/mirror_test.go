package envset

import (
	"testing"
)

func baseMirrorSets(t *testing.T) (src, dst *EnvSet) {
	t.Helper()
	src, _ = New("src", "local")
	src.Vars["API_KEY"] = "abc123"
	src.Vars["DB_HOST"] = "localhost"
	src.Vars["SECRET"] = "topsecret"

	dst, _ = New("dst", "local")
	dst.Vars["EXISTING"] = "keep"
	return
}

func TestMirror_AllKeys(t *testing.T) {
	src, dst := baseMirrorSets(t)
	n, err := Mirror(src, dst, DefaultMirrorOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 keys mirrored, got %d", n)
	}
	if dst.Vars["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123")
	}
	if dst.Vars["EXISTING"] != "keep" {
		t.Errorf("EXISTING should be untouched")
	}
}

func TestMirror_WithPrefix(t *testing.T) {
	src, dst := baseMirrorSets(t)
	opts := DefaultMirrorOptions()
	opts.Prefix = "MIRROR_"
	n, err := Mirror(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 mirrored, got %d", n)
	}
	if dst.Vars["MIRROR_API_KEY"] != "abc123" {
		t.Errorf("expected MIRROR_API_KEY")
	}
}

func TestMirror_NoOverwrite(t *testing.T) {
	src, dst := baseMirrorSets(t)
	dst.Vars["API_KEY"] = "original"
	n, err := Mirror(src, dst, DefaultMirrorOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 mirrored (skipped existing), got %d", n)
	}
	if dst.Vars["API_KEY"] != "original" {
		t.Errorf("API_KEY should not be overwritten")
	}
}

func TestMirror_WithOverwrite(t *testing.T) {
	src, dst := baseMirrorSets(t)
	dst.Vars["API_KEY"] = "original"
	opts := DefaultMirrorOptions()
	opts.Overwrite = true
	_, err := Mirror(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Vars["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY to be overwritten")
	}
}

func TestMirror_SelectedKeys(t *testing.T) {
	src, dst := baseMirrorSets(t)
	opts := DefaultMirrorOptions()
	opts.Keys = []string{"API_KEY"}
	n, err := Mirror(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 mirrored, got %d", n)
	}
	if _, ok := dst.Vars["DB_HOST"]; ok {
		t.Errorf("DB_HOST should not be mirrored")
	}
}

func TestMirror_MissingKey(t *testing.T) {
	src, dst := baseMirrorSets(t)
	opts := DefaultMirrorOptions()
	opts.Keys = []string{"NONEXISTENT"}
	_, err := Mirror(src, dst, opts)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestMirror_NilSource(t *testing.T) {
	_, dst := baseMirrorSets(t)
	_, err := Mirror(nil, dst, DefaultMirrorOptions())
	if err == nil {
		t.Error("expected error for nil source")
	}
}
