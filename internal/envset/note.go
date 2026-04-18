package envset

import (
	"fmt"
	"strings"
	"time"
)

// Note represents a free-form note attached to an EnvSet.
type Note struct {
	Text      string
	CreatedAt time.Time
	Author    string
}

// AddNote attaches a note to the EnvSet metadata.
func AddNote(es *EnvSet, text, author string) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("note text must not be empty")
	}
	if strings.ContainsAny(text, "\n\r") {
		return fmt.Errorf("note text must not contain newlines")
	}
	n := Note{Text: text, Author: author, CreatedAt: time.Now().UTC()}
	key := fmt.Sprintf("__note_%d", len(ListNotes(es)))
	es.Meta[key] = fmt.Sprintf("%s|%s|%s", n.CreatedAt.Format(time.RFC3339), n.Author, n.Text)
	return nil
}

// ListNotes returns all notes attached to the EnvSet.
func ListNotes(es *EnvSet) []Note {
	if es == nil {
		return nil
	}
	var notes []Note
	for i := 0; ; i++ {
		key := fmt.Sprintf("__note_%d", i)
		val, ok := es.Meta[key]
		if !ok {
			break
		}
		parts := strings.SplitN(val, "|", 3)
		if len(parts) != 3 {
			continue
		}
		t, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		notes = append(notes, Note{CreatedAt: t, Author: parts[1], Text: parts[2]})
	}
	return notes
}

// ClearNotes removes all notes from the EnvSet.
func ClearNotes(es *EnvSet) error {
	if es == nil {
		return fmt.Errorf("envset is nil")
	}
	for i := 0; ; i++ {
		key := fmt.Sprintf("__note_%d", i)
		if _, ok := es.Meta[key]; !ok {
			break
		}
		delete(es.Meta, key)
	}
	return nil
}
