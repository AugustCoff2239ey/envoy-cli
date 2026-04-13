package envset

import (
	"testing"
)

func baseTemplateSet() *EnvSet {
	e, _ := New("app", "local")
	e.Vars["DB_HOST"] = "localhost"
	e.Vars["DB_PORT"] = "5432"
	e.Vars["APP_NAME"] = "envoy"
	return e
}

func TestRenderTemplate_AllResolved(t *testing.T) {
	e := baseTemplateSet()
	res, err := RenderTemplate(e, "host={{DB_HOST}} port={{DB_PORT}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "host=localhost port=5432"
	if res.Rendered != want {
		t.Errorf("got %q, want %q", res.Rendered, want)
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestRenderTemplate_Unresolved(t *testing.T) {
	e := baseTemplateSet()
	res, err := RenderTemplate(e, "host={{DB_HOST}} user={{DB_USER}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "DB_USER" {
		t.Errorf("expected [DB_USER] unresolved, got %v", res.Unresolved)
	}
	if res.Rendered != "host=localhost user={{DB_USER}}" {
		t.Errorf("unexpected rendered: %q", res.Rendered)
	}
}

func TestRenderTemplate_EmptyTemplate(t *testing.T) {
	e := baseTemplateSet()
	res, err := RenderTemplate(e, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "" {
		t.Errorf("expected empty rendered, got %q", res.Rendered)
	}
}

func TestRenderTemplate_NilEnvSet(t *testing.T) {
	_, err := RenderTemplate(nil, "{{DB_HOST}}")
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestExtractPlaceholders_Unique(t *testing.T) {
	keys := ExtractPlaceholders("{{A}} {{B}} {{A}} {{C}}")
	if len(keys) != 3 {
		t.Errorf("expected 3 unique keys, got %d: %v", len(keys), keys)
	}
}

func TestTemplateComplete_True(t *testing.T) {
	e := baseTemplateSet()
	if !TemplateComplete(e, "{{DB_HOST}}:{{DB_PORT}}") {
		t.Error("expected template to be complete")
	}
}

func TestTemplateComplete_False(t *testing.T) {
	e := baseTemplateSet()
	if TemplateComplete(e, "{{DB_HOST}}:{{MISSING}}") {
		t.Error("expected template to be incomplete")
	}
}

func TestMissingPlaceholders(t *testing.T) {
	e := baseTemplateSet()
	missing := MissingPlaceholders(e, "{{DB_HOST}} {{FOO}} {{BAR}}")
	if len(missing) != 2 {
		t.Errorf("expected 2 missing, got %v", missing)
	}
}

func TestSuggestPlaceholders(t *testing.T) {
	e := baseTemplateSet()
	sugs := SuggestPlaceholders(e, []string{"DB_HOS"})
	if len(sugs["DB_HOS"]) == 0 {
		t.Error("expected suggestion for DB_HOS -> DB_HOST")
	}
}
