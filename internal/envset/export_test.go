package envset

import (
	"strings"
	"testing"
)

func baseExportSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("test", "local")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("APP_ENV", "local")
	_ = es.Set("DB_URL", "postgres://localhost/dev")
	_ = es.Set("PORT", "8080")
	return es
}

func TestExport_Dotenv(t *testing.T) {
	es := baseExportSet(t)
	out, err := Export(es, FormatDotenv)
	if err != nil {
		t.Fatalf("Export dotenv: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=local") {
		t.Errorf("expected APP_ENV=local in dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in dotenv output, got:\n%s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	es := baseExportSet(t)
	out, err := Export(es, FormatExport)
	if err != nil {
		t.Fatalf("Export shell: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected 'export APP_ENV=' in shell output, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	es := baseExportSet(t)
	out, err := Export(es, FormatJSON)
	if err != nil {
		t.Fatalf("Export json: %v", err)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected \"APP_ENV\" in json output, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected valid JSON object braces, got:\n%s", out)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	es := baseExportSet(t)
	_, err := Export(es, Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyEnvSet(t *testing.T) {
	es, _ := New("empty", "local")
	out, err := Export(es, FormatDotenv)
	if err != nil {
		t.Fatalf("Export empty: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output for empty EnvSet, got: %q", out)
	}
}
