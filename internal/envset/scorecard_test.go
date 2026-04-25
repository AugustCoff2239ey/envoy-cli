package envset

import (
	"testing"
)

func baseScorecardSet(t *testing.T) *EnvSet {
	t.Helper()
	es, err := New("scorecard-set", "staging")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = es.Set("DB_HOST", "localhost")
	_ = es.Set("DB_PORT", "5432")
	_ = es.Set("APP_ENV", "staging")
	return es
}

func TestScorecard_BasicScore(t *testing.T) {
	es := baseScorecardSet(t)
	res, err := Scorecard(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score < 0 || res.Score > 100 {
		t.Errorf("score %d out of range [0,100]", res.Score)
	}
	if res.Grade == "" {
		t.Error("grade should not be empty")
	}
}

func TestScorecard_NamingPenalty(t *testing.T) {
	es, _ := New("bad-names", "local")
	_ = es.Set("bad-key", "value")
	_ = es.Set("another-bad", "value")
	res, err := Scorecard(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Breakdown["naming"] > 5 {
		t.Errorf("expected low naming score for bad keys, got %d", res.Breakdown["naming"])
	}
	foundSuggestion := false
	for _, s := range res.Suggestions {
		if len(s) > 0 {
			foundSuggestion = true
		}
	}
	if !foundSuggestion {
		t.Error("expected at least one suggestion")
	}
}

func TestScorecard_SecurityBonus(t *testing.T) {
	es := baseScorecardSet(t)
	_ = es.Set("DB_PASSWORD", "s3cr3t")
	_ = LockKey(es, "DB_PASSWORD", "system")
	res, err := Scorecard(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Breakdown["security"] < 20 {
		t.Errorf("expected higher security score when sensitive key is locked, got %d", res.Breakdown["security"])
	}
}

func TestScorecard_EmptySet(t *testing.T) {
	es, _ := New("empty", "local")
	res, err := Scorecard(es)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 0 {
		t.Errorf("expected score 0 for empty set, got %d", res.Score)
	}
	if res.Grade != "F" {
		t.Errorf("expected grade F for empty set, got %s", res.Grade)
	}
}

func TestScorecard_NilEnvSet(t *testing.T) {
	_, err := Scorecard(nil)
	if err == nil {
		t.Error("expected error for nil EnvSet")
	}
}

func TestScorecard_GradeThresholds(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{95, "A"}, {80, "B"}, {65, "C"}, {45, "D"}, {20, "F"},
	}
	for _, tc := range cases {
		got := scoreGrade(tc.score)
		if got != tc.want {
			t.Errorf("scoreGrade(%d) = %q, want %q", tc.score, got, tc.want)
		}
	}
}
