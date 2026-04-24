package envset

import (
	"testing"
)

func baseRevealSet() *EnvSet {
	es, _ := New("reveal-test", "staging")
	_ = es.Set("API_KEY", "supersecretvalue")
	_ = es.Set("DB_PASSWORD", "p@ssw0rd!")
	_ = es.Set("SHORT", "ab")
	_ = es.Set("PLAIN", "hello")
	return es
}

func TestReveal_DefaultOptions(t *testing.T) {
	es := baseRevealSet()
	opts := DefaultRevealOptions()
	results, err := Reveal(es, []string{"API_KEY"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// "supersecretvalue" has 16 chars; suffix=4 => "alue" visible
	if results[0].Visible != "************alue" {
		t.Errorf("unexpected visible: %q", results[0].Visible)
	}
}

func TestReveal_WithPrefixAndSuffix(t *testing.T) {
	es := baseRevealSet()
	opts := RevealOptions{PrefixLen: 2, SuffixLen: 2}
	results, err := Reveal(es, []string{"DB_PASSWORD"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "p@ssw0rd!" has 9 chars; prefix=2 => "p@", suffix=2 => "d!", middle=5 => "*****"
	expected := "p@*****d!"
	if results[0].Visible != expected {
		t.Errorf("expected %q, got %q", expected, results[0].Visible)
	}
}

func TestReveal_ShortValue(t *testing.T) {
	es := baseRevealSet()
	opts := RevealOptions{PrefixLen: 1, SuffixLen: 4}
	results, err := Reveal(es, []string{"SHORT"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// shown >= len, so full value returned
	if results[0].Visible != "ab" {
		t.Errorf("expected full value %q, got %q", "ab", results[0].Visible)
	}
}

func TestReveal_MissingKey(t *testing.T) {
	es := baseRevealSet()
	opts := DefaultRevealOptions()
	_, err := Reveal(es, []string{"MISSING"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestReveal_NilEnvSet(t *testing.T) {
	opts := DefaultRevealOptions()
	_, err := Reveal(nil, []string{"API_KEY"}, opts)
	if err != ErrNilEnvSet {
		t.Errorf("expected ErrNilEnvSet, got %v", err)
	}
}

func TestReveal_NegativeSuffixLen(t *testing.T) {
	es := baseRevealSet()
	opts := RevealOptions{SuffixLen: -1}
	_, err := Reveal(es, []string{"PLAIN"}, opts)
	if err == nil {
		t.Fatal("expected error for negative SuffixLen")
	}
}

func TestReveal_MultipleKeys(t *testing.T) {
	es := baseRevealSet()
	opts := DefaultRevealOptions()
	results, err := Reveal(es, []string{"API_KEY", "PLAIN"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
