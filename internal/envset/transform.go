package envset

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms an env var value.
type TransformFunc func(value string) (string, error)

// TransformOptions controls how Transform behaves.
type TransformOptions struct {
	// Keys limits transformation to specific keys; empty means all keys.
	Keys []string
	// SkipLocked skips locked keys instead of returning an error.
	SkipLocked bool
}

// builtinTransforms maps transform names to their implementations.
var builtinTransforms = map[string]TransformFunc{
	"upper":   func(v string) (string, error) { return strings.ToUpper(v), nil },
	"lower":   func(v string) (string, error) { return strings.ToLower(v), nil },
	"trim":    func(v string) (string, error) { return strings.TrimSpace(v), nil },
	"reverse": func(v string) (string, error) { return reverseString(v), nil },
}

// Transform applies a named or custom transform function to env var values.
// If opts.Keys is empty, all keys are transformed.
func Transform(es *EnvSet, transformName string, fn TransformFunc, opts TransformOptions) error {
	if es == nil {
		return fmt.Errorf("transform: envset is nil")
	}

	if fn == nil {
		builtin, ok := builtinTransforms[transformName]
		if !ok {
			return fmt.Errorf("transform: unknown transform %q", transformName)
		}
		fn = builtin
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range es.Vars {
			targetKeys = append(targetKeys, k)
		}
	}

	for _, k := range targetKeys {
		val, exists := es.Vars[k]
		if !exists {
			return fmt.Errorf("transform: key %q not found", k)
		}
		if IsLocked(es, k) {
			if opts.SkipLocked {
				continue
			}
			return fmt.Errorf("transform: key %q is locked", k)
		}
		newVal, err := fn(val)
		if err != nil {
			return fmt.Errorf("transform: key %q: %w", k, err)
		}
		es.Vars[k] = newVal
	}
	return nil
}

func reverseString(s string) (string, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}
