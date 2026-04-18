package envset

import (
	"fmt"
	"strings"
)

// ClassifyResult holds the classification of a single key.
type ClassifyResult struct {
	Key      string
	Category string
}

// ClassifyReport holds all classification results for an EnvSet.
type ClassifyReport struct {
	Results []ClassifyResult
	ByCategory map[string][]string
}

var builtinCategories = map[string][]string{
	"secret":   {"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH"},
	"database": {"DB", "DATABASE", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN"},
	"network":  {"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "DOMAIN"},
	"feature":  {"FEATURE", "FLAG", "ENABLE", "DISABLE", "TOGGLE"},
}

// Classify categorises each key in the EnvSet based on keyword matching.
func Classify(es *EnvSet) (*ClassifyReport, error) {
	if es == nil {
		return nil, fmt.Errorf("classify: nil EnvSet")
	}

	report := &ClassifyReport{
		ByCategory: make(map[string][]string),
	}

	for k := range es.Vars {
		upper := strings.ToUpper(k)
		category := "general"

		outer:
		for cat, keywords := range builtinCategories {
			for _, kw := range keywords {
				if strings.Contains(upper, kw) {
					category = cat
					break outer
				}
			}
		}

		report.Results = append(report.Results, ClassifyResult{Key: k, Category: category})
		report.ByCategory[category] = append(report.ByCategory[category], k)
	}

	return report, nil
}
