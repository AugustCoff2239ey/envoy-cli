package envset

import (
	"fmt"
	"strings"
)

// ScorecardResult holds the quality score and breakdown for an EnvSet.
type ScorecardResult struct {
	Score      int               // 0-100
	Grade      string            // A, B, C, D, F
	Breakdown  map[string]int    // category -> points awarded
	Suggestions []string
}

// Scorecard evaluates the quality of an EnvSet and returns a ScorecardResult.
// It checks documentation, key naming, value hygiene, and security signals.
func Scorecard(es *EnvSet) (*ScorecardResult, error) {
	if es == nil {
		return nil, fmt.Errorf("scorecard: nil EnvSet")
	}

	breakdown := map[string]int{
		"naming":        0,
		"documentation": 0,
		"hygiene":       0,
		"security":      0,
	}
	var suggestions []string

	keys := es.Keys()
	if len(keys) == 0 {
		return &ScorecardResult{Score: 0, Grade: "F", Breakdown: breakdown,
			Suggestions: []string{"add environment variables to the set"}}, nil
	}

	// Naming: all keys uppercase and underscore-only separators
	namingOK := 0
	for _, k := range keys {
		if k == strings.ToUpper(k) && !strings.ContainsAny(k, "-. ") {
			namingOK++
		}
	}
	naming := namingOK * 25 / len(keys)
	breakdown["naming"] = naming
	if naming < 20 {
		suggestions = append(suggestions, "use UPPER_SNAKE_CASE for all keys")
	}

	// Documentation: keys with annotations
	annotated := 0
	for _, k := range keys {
		if _, err := GetAnnotation(es, k); err == nil {
			annotated++
		}
	}
	doc := annotated * 25 / len(keys)
	breakdown["documentation"] = doc
	if doc < 15 {
		suggestions = append(suggestions, "annotate keys with descriptions using 'annotate'")
	}

	// Hygiene: no empty values
	nonEmpty := 0
	for _, k := range keys {
		if v, _ := es.Get(k); strings.TrimSpace(v) != "" {
			nonEmpty++
		}
	}
	hyg := nonEmpty * 25 / len(keys)
	breakdown["hygiene"] = hyg
	if hyg < 20 {
		suggestions = append(suggestions, "remove or populate empty-valued keys")
	}

	// Security: sensitive keys are locked or protected
	sensitiveTotal := 0
	sensitiveSecured := 0
	for _, k := range keys {
		if IsSensitiveKey(k) {
			sensitiveTotal++
			if IsLocked(es, k) || IsProtected(es, k) {
				sensitiveSecured++
			}
		}
	}
	sec := 25
	if sensitiveTotal > 0 {
		sec = sensitiveSecured * 25 / sensitiveTotal
		if sec < 20 {
			suggestions = append(suggestions, "lock or protect sensitive keys (SECRET, PASSWORD, TOKEN, etc.)")
		}
	}
	breakdown["security"] = sec

	total := naming + doc + hyg + sec
	grade := scoreGrade(total)

	return &ScorecardResult{
		Score:       total,
		Grade:       grade,
		Breakdown:   breakdown,
		Suggestions: suggestions,
	}, nil
}

func scoreGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
