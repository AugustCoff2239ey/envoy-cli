package envset

import (
	"testing"
)

func baseBookmarkSet() *EnvSet {
	e, _ := New("myapp", "staging")
	_ = e.Set("DB_URL", "postgres://localhost/db")
	_ = e.Set("API_KEY", "secret")
	return e
}

func TestAddBookmark_Valid(t *testing.T) {
	e := baseBookmarkSet()
	if err := AddBookmark(e, "my-bookmark"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasBookmark(e, "my-bookmark") {
		t.Error("expected bookmark to exist")
	}
}

func TestAddBookmark_EmptyName(t *testing.T) {
	e := baseBookmarkSet()
	if err := AddBookmark(e, ""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestAddBookmark_InvalidChars(t *testing.T) {
	e := baseBookmarkSet()
	if err := AddBookmark(e, "bad name!"); err == nil {
		t.Error("expected error for invalid characters")
	}
}

func TestAddBookmark_NilEnvSet(t *testing.T) {
	if err := AddBookmark(nil, "bm"); err == nil {
		t.Error("expected error for nil envset")
	}
}

func TestRemoveBookmark_Valid(t *testing.T) {
	e := baseBookmarkSet()
	_ = AddBookmark(e, "to-remove")
	if err := RemoveBookmark(e, "to-remove"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasBookmark(e, "to-remove") {
		t.Error("expected bookmark to be removed")
	}
}

func TestRemoveBookmark_NotFound(t *testing.T) {
	e := baseBookmarkSet()
	if err := RemoveBookmark(e, "ghost"); err == nil {
		t.Error("expected error for missing bookmark")
	}
}

func TestListBookmarks_Multiple(t *testing.T) {
	e := baseBookmarkSet()
	_ = AddBookmark(e, "bm1")
	_ = AddBookmark(e, "bm2")
	bms := ListBookmarks(e)
	if len(bms) != 2 {
		t.Fatalf("expected 2 bookmarks, got %d", len(bms))
	}
}

func TestListBookmarks_NilEnvSet(t *testing.T) {
	if bms := ListBookmarks(nil); bms != nil {
		t.Error("expected nil for nil envset")
	}
}

func TestBookmark_StoresSetAndEnv(t *testing.T) {
	e := baseBookmarkSet()
	_ = AddBookmark(e, "ref")
	bms := ListBookmarks(e)
	if len(bms) != 1 {
		t.Fatalf("expected 1 bookmark")
	}
	if bms[0].SetName != "myapp" || bms[0].Env != "staging" {
		t.Errorf("unexpected bookmark data: %+v", bms[0])
	}
}
