package envset

import (
	"fmt"
	"strings"
)

// FilterOptions controls how keys are selected during filtering.
type FilterOptions struct {
	Prefix    string
	Suffix    string
	Envs      []string
	Exclude   []string
}

// Filter returns a new EnvSet containing only keys matching the given options.
func Filter(es *EnvSet, opts FilterOptions) (*EnvSet, error) {
	if es == nil {
		return nil, fmt.Errorf("filter: nil EnvSet")
	}

	excludeSet := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[strings.ToUpper(k)] = true
	}

	envSet := make(map[string]bool, len(opts.Envs))
	for _, e := range opts.Envs {
		envSet[strings.ToLower(e)] = true
	}

	out, err := New(es.Name, es.Environment)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	if len(envSet) > 0 && !envSet[strings.ToLower(es.Environment)] {
		return out, nil
	}

	for k, v := range es.Vars {
		if excludeSet[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, strings.ToUpper(opts.Prefix)) {
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, strings.ToUpper(opts.Suffix)) {
			continue
		}
		if err := out.Set(k, v); err != nil {
			return nil, fmt.Errorf("filter: set %s: %w", k, err)
		}
	}

	return out, nil
}
