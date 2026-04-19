package envset

import "time"

// StatusReport summarizes the live state of an EnvSet.
type StatusReport struct {
	Name        string
	Environment string
	TotalKeys   int
	LockedKeys  int
	PinnedKeys  int
	Protected   int
	ExpiredKeys int
	DraftKeys   int
	Readonly    bool
	Frozen      bool
	GeneratedAt time.Time
}

// Status returns a StatusReport for the given EnvSet.
func Status(es *EnvSet) (StatusReport, error) {
	if es == nil {
		return StatusReport{}, ErrNilEnvSet
	}

	report := StatusReport{
		Name:        es.Name,
		Environment: es.Environment,
		TotalKeys:   len(es.Vars),
		Readonly:    IsReadonly(es),
		Frozen:      IsFrozen(es),
		GeneratedAt: time.Now().UTC(),
	}

	for k := range es.Vars {
		if IsLocked(es, k) {
			report.LockedKeys++
		}
		if IsPinned(es, k) {
			report.PinnedKeys++
		}
		if IsProtected(es, k) {
			report.Protected++
		}
		if IsExpired(es, k) {
			report.ExpiredKeys++
		}
	}

	if v, ok := es.Meta["draft"]; ok && v == "true" {
		report.DraftKeys = report.TotalKeys
	}

	return report, nil
}
