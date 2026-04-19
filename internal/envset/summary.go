package envset

import (
	"fmt"
	"sort"
	"strings"
)

// SummaryReport holds a high-level overview of an EnvSet.
type SummaryReport struct {
	Name        string
	Environment string
	TotalKeys   int
	LockedKeys  []string
	Protected   []string
	Pinned      []string
	TagList     []string
	Expired     []string
	Readonly    bool
	Frozen      bool
}

// Summary generates a SummaryReport for the given EnvSet.
func Summary(es *EnvSet) (*SummaryReport, error) {
	if es == nil {
		return nil, fmt.Errorf("summary: envset is nil")
	}

	r := &SummaryReport{
		Name:        es.Name,
		Environment: es.Environment,
		TotalKeys:   len(es.Vars),
		Readonly:    IsReadonly(es),
		Frozen:      IsFrozen(es),
	}

	for k := range es.Vars {
		if IsLocked(es, k) {
			r.LockedKeys = append(r.LockedKeys, k)
		}
		if IsProtected(es, k) {
			r.Protected = append(r.Protected, k)
		}
		if IsPinned(es, k) {
			r.Pinned = append(r.Pinned, k)
		}
		if IsExpired(es, k) {
			r.Expired = append(r.Expired, k)
		}
	}

	if tags, ok := es.Meta["tags"]; ok && tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				r.TagList = append(r.TagList, t)
			}
		}
	}

	sort.Strings(r.LockedKeys)
	sort.Strings(r.Protected)
	sort.Strings(r.Pinned)
	sort.Strings(r.Expired)

	return r, nil
}
