package envset

import (
	"testing"
)

func baseDraftSet() *EnvSet {
	e, _ := New("myapp", "staging")
	_ = e.Set("API_URL", "https://staging.example.com")
	_ = e.Set("DB_PASS", "secret")
	return e
}

func TestSaveDraft_Basic(t *testing.T) {
	e := baseDraftSet()
	d, err := SaveDraft(e, "wip changes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "draft:myapp" {
		t.Errorf("expected draft:myapp, got %s", d.Name)
	}
	if d.Meta["draft_note"] != "wip changes" {
		t.Errorf("expected note 'wip changes', got %s", d.Meta["draft_note"])
	}
}

func TestSaveDraft_Independence(t *testing.T) {
	e := baseDraftSet()
	d, _ := SaveDraft(e, "")
	d.Vars["API_URL"] = "changed"
	if e.Vars["API_URL"] == "changed" {
		t.Error("draft mutation should not affect source")
	}
}

func TestSaveDraft_NilSource(t *testing.T) {
	_, err := SaveDraft(nil, "note")
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestPromoteDraft_Valid(t *testing.T) {
	e := baseDraftSet()
	d, _ := SaveDraft(e, "note")
	p, err := PromoteDraft(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "myapp" {
		t.Errorf("expected myapp, got %s", p.Name)
	}
	if _, ok := p.Meta["draft_note"]; ok {
		t.Error("draft_note should be removed after promotion")
	}
}

func TestPromoteDraft_NotADraft(t *testing.T) {
	e := baseDraftSet()
	_, err := PromoteDraft(e)
	if err == nil {
		t.Error("expected error promoting non-draft")
	}
}

func TestPromoteDraft_Nil(t *testing.T) {
	_, err := PromoteDraft(nil)
	if err == nil {
		t.Error("expected error for nil draft")
	}
}

func TestIsDraft(t *testing.T) {
	e := baseDraftSet()
	if IsDraft(e) {
		t.Error("regular envset should not be draft")
	}
	d, _ := SaveDraft(e, "")
	if !IsDraft(d) {
		t.Error("saved draft should be detected as draft")
	}
	if IsDraft(nil) {
		t.Error("nil should not be draft")
	}
}
