package envset

import "fmt"

// QuotaOptions defines limits for an EnvSet.
type QuotaOptions struct {
	MaxKeys   int
	MaxKeyLen int
	MaxValLen int
}

// DefaultQuotaOptions returns sensible defaults.
func DefaultQuotaOptions() QuotaOptions {
	return QuotaOptions{
		MaxKeys:   100,
		MaxKeyLen: 64,
		MaxValLen: 1024,
	}
}

// QuotaViolation describes a single quota breach.
type QuotaViolation struct {
	Key     string
	Message string
}

// CheckQuota validates an EnvSet against the given QuotaOptions.
// It returns a slice of violations; an empty slice means the set is within quota.
func CheckQuota(es *EnvSet, opts QuotaOptions) ([]QuotaViolation, error) {
	if es == nil {
		return nil, fmt.Errorf("envset: CheckQuota called on nil EnvSet")
	}

	var violations []QuotaViolation

	if opts.MaxKeys > 0 && len(es.Vars) > opts.MaxKeys {
		violations = append(violations, QuotaViolation{
			Key:     "__count__",
			Message: fmt.Sprintf("key count %d exceeds max %d", len(es.Vars), opts.MaxKeys),
		})
	}

	for k, v := range es.Vars {
		if opts.MaxKeyLen > 0 && len(k) > opts.MaxKeyLen {
			violations = append(violations, QuotaViolation{
				Key:     k,
				Message: fmt.Sprintf("key length %d exceeds max %d", len(k), opts.MaxKeyLen),
			})
		}
		if opts.MaxValLen > 0 && len(v) > opts.MaxValLen {
			violations = append(violations, QuotaViolation{
				Key:     k,
				Message: fmt.Sprintf("value length %d exceeds max %d", len(v), opts.MaxValLen),
			})
		}
	}

	return violations, nil
}
