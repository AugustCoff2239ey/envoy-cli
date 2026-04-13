package envset

import (
	"testing"
)

func baseCopySets(t *testing.T) (*EnvSet, *EnvSet) {
	t.Helper()
	src, _ := New("src", "staging")
	src.Vars["FOO"] = "foo_val"
	src.Vars["BAR"] = "bar_val"
	src.Vars["BAZ"] = "baz_val"

	dst, _ := New("dst", "staging")
	dst.Vars["BAR"] = "existing_bar"
	return src, dst
}

func TestCopy_AllKeys(t *testing.T) {
	src, dst := baseCopySets(t)
	n, err := Copy(src, dst, CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 copied, got %d", n)
	}
	if dst.Vars["FOO"] != "foo_val" {
		t.Errorf("expected FOO=foo_val, got %q", dst.Vars["FOO"])
	}
	if dst.Vars["BAR"] != "bar_val" {
		t.Errorf("expected BAR=bar_val (overwritten), got %q", dst.Vars["BAR"])
	}
}

func TestCopy_NoOverwrite(t *testing.T) {
	src, dst := baseCopySets(t)
	n, err := Copy(src, dst, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// BAR already exists and Overwrite=false, so only FOO and BAZ are copied
	if n != 2 {
		t.Errorf("expected 2 copied, got %d", n)
	}
	if dst.Vars["BAR"] != "existing_bar" {
		t.Errorf("expected BAR to remain existing_bar, got %q", dst.Vars["BAR"])
	}
}

func TestCopy_SelectedKeys(t *testing.T) {
	src, dst := baseCopySets(t)
	n, err := Copy(src, dst, CopyOptions{Overwrite: true, Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 copied, got %d", n)
	}
	if _, ok := dst.Vars["BAZ"]; ok {
		t.Error("BAZ should not have been copied")
	}
}

func TestCopy_MissingKey(t *testing.T) {
	src, dst := baseCopySets(t)
	_, err := Copy(src, dst, CopyOptions{Keys: []string{"MISSING"}})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestCopy_NilSource(t *testing.T) {
	dst, _ := New("dst", "local")
	_, err := Copy(nil, dst, CopyOptions{})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestCopy_NilDestination(t *testing.T) {
	src, _ := New("src", "local")
	_, err := Copy(src, nil, CopyOptions{})
	if err == nil {
		t.Fatal("expected error for nil destination")
	}
}
