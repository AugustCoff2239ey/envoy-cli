package envset

import "fmt"

// DraftMeta holds metadata about a draft envset.
type DraftMeta struct {
	Name        string
	Environment string
	Note        string
}

// SaveDraft stores a copy of the envset as a named draft without persisting
// it to the main store. Drafts are kept in a separate namespace.
func SaveDraft(e *EnvSet, note string) (*EnvSet, error) {
	if e == nil {
		return nil, fmt.Errorf("draft: source envset is nil")
	}
	draft := &EnvSet{
		Name:        "draft:" + e.Name,
		Environment: e.Environment,
		Vars:        make(map[string]string, len(e.Vars)),
	}
	for k, v := range e.Vars {
		draft.Vars[k] = v
	}
	if draft.Meta == nil {
		draft.Meta = make(map[string]string)
	}
	for k, v := range e.Meta {
		draft.Meta[k] = v
	}
	draft.Meta["draft_note"] = note
	return draft, nil
}

// PromoteDraft converts a draft envset back to a regular envset, stripping the
// draft prefix and note.
func PromoteDraft(draft *EnvSet) (*EnvSet, error) {
	if draft == nil {
		return nil, fmt.Errorf("draft: draft envset is nil")
	}
	const prefix = "draft:"
	if len(draft.Name) <= len(prefix) || draft.Name[:len(prefix)] != prefix {
		return nil, fmt.Errorf("draft: envset %q is not a draft", draft.Name)
	}
	promoted := &EnvSet{
		Name:        draft.Name[len(prefix):],
		Environment: draft.Environment,
		Vars:        make(map[string]string, len(draft.Vars)),
	}
	for k, v := range draft.Vars {
		promoted.Vars[k] = v
	}
	promoted.Meta = make(map[string]string)
	for k, v := range draft.Meta {
		promoted.Meta[k] = v
	}
	delete(promoted.Meta, "draft_note")
	return promoted, nil
}

// IsDraft reports whether the envset is a draft.
func IsDraft(e *EnvSet) bool {
	if e == nil {
		return false
	}
	return len(e.Name) > 6 && e.Name[:6] == "draft:"
}
