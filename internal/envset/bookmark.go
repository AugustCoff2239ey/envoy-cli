package envset

import (
	"fmt"
	"regexp"
)

var validBookmarkName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Bookmark represents a named pointer to an envset.
type Bookmark struct {
	Name    string `json:"name"`
	SetName string `json:"set_name"`
	Env     string `json:"env"`
}

// AddBookmark adds a named bookmark referencing the given envset name and environment.
func AddBookmark(e *EnvSet, name string) error {
	if e == nil {
		return fmt.Errorf("envset is nil")
	}
	if name == "" {
		return fmt.Errorf("bookmark name must not be empty")
	}
	if !validBookmarkName.MatchString(name) {
		return fmt.Errorf("bookmark name %q contains invalid characters", name)
	}
	if e.Meta == nil {
		e.Meta = map[string]string{}
	}
	key := "bookmark:" + name
	e.Meta[key] = e.Name + "|" + e.Environment
	return nil
}

// RemoveBookmark removes a named bookmark from the envset.
func RemoveBookmark(e *EnvSet, name string) error {
	if e == nil {
		return fmt.Errorf("envset is nil")
	}
	key := "bookmark:" + name
	if _, ok := e.Meta[key]; !ok {
		return fmt.Errorf("bookmark %q not found", name)
	}
	delete(e.Meta, key)
	return nil
}

// ListBookmarks returns all bookmarks stored on the envset.
func ListBookmarks(e *EnvSet) []Bookmark {
	if e == nil {
		return nil
	}
	var out []Bookmark
	for k, v := range e.Meta {
		if len(k) > 9 && k[:9] == "bookmark:" {
			bname := k[9:]
			setName, env := splitPipe(v)
			out = append(out, Bookmark{Name: bname, SetName: setName, Env: env})
		}
	}
	return out
}

// HasBookmark returns true if the named bookmark exists.
func HasBookmark(e *EnvSet, name string) bool {
	if e == nil {
		return false
	}
	_, ok := e.Meta["bookmark:"+name]
	return ok
}

func splitPipe(s string) (string, string) {
	for i, c := range s {
		if c == '|' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}
