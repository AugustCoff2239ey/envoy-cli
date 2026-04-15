package envset

import (
	"fmt"
	"strings"
	"unicode"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all key names to UPPER_SNAKE_CASE.
	UppercaseKeys bool

	// TrimValues removes leading and trailing whitespace from values.
	TrimValues bool

	// ReplaceHyphens replaces hyphens in keys with underscores.
	ReplaceHyphens bool

	// StripNonPrintable removes non-printable characters from values.
	StripNonPrintable bool

	// Keys restricts normalization to the specified keys.
	// If empty, all keys are normalized.
	Keys []string
}

// DefaultNormalizeOptions returns a NormalizeOptions with sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys:     true,
		TrimValues:        true,
		ReplaceHyphens:    true,
		StripNonPrintable: false,
	}
}

// NormalizeResult holds the outcome of a Normalize operation.
type NormalizeResult struct {
	// KeysRenamed maps original key names to their normalized forms.
	KeysRenamed map[string]string

	// ValuesChanged lists keys whose values were modified.
	ValuesChanged []string
}

// Normalize applies normalization rules to the keys and values of an EnvSet
// according to the provided options. It returns a NormalizeResult describing
// what changed, or an error if a normalized key would conflict with an
// existing key or violate validation rules.
func Normalize(es *EnvSet, opts NormalizeOptions) (NormalizeResult, error) {
	if es == nil {
		return NormalizeResult{}, fmt.Errorf("normalize: envset must not be nil")
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range es.Vars {
			targetKeys = append(targetKeys, k)
		}
	}

	result := NormalizeResult{
		KeysRenamed: make(map[string]string),
	}

	// Collect renames first to detect conflicts before mutating.
	renames := make(map[string]string) // old -> new
	for _, key := range targetKeys {
		newKey := normalizeKey(key, opts)
		if newKey != key {
			// Ensure the new key doesn't already exist (and isn't itself being renamed).
			if _, exists := es.Vars[newKey]; exists {
				return NormalizeResult{}, fmt.Errorf(
					"normalize: key %q would conflict with existing key %q after normalization",
					key, newKey,
				)
			}
			if err := ValidateKey(newKey); err != nil {
				return NormalizeResult{}, fmt.Errorf("normalize: normalized key %q is invalid: %w", newKey, err)
			}
			renames[key] = newKey
		}
	}

	// Apply key renames.
	for oldKey, newKey := range renames {
		es.Vars[newKey] = es.Vars[oldKey]
		delete(es.Vars, oldKey)
		result.KeysRenamed[oldKey] = newKey
	}

	// Apply value normalization. Use the new key name if a rename occurred.
	for _, key := range targetKeys {
		actualKey := key
		if renamed, ok := renames[key]; ok {
			actualKey = renamed
		}
		val, exists := es.Vars[actualKey]
		if !exists {
			continue
		}
		newVal := normalizeValue(val, opts)
		if newVal != val {
			es.Vars[actualKey] = newVal
			result.ValuesChanged = append(result.ValuesChanged, actualKey)
		}
	}

	return result, nil
}

// normalizeKey applies key-level transformations based on opts.
func normalizeKey(key string, opts NormalizeOptions) string {
	if opts.ReplaceHyphens {
		key = strings.ReplaceAll(key, "-", "_")
	}
	if opts.UppercaseKeys {
		key = strings.ToUpper(key)
	}
	return key
}

// normalizeValue applies value-level transformations based on opts.
func normalizeValue(val string, opts NormalizeOptions) string {
	if opts.TrimValues {
		val = strings.TrimSpace(val)
	}
	if opts.StripNonPrintable {
		val = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) || r == '\t' {
				return r
			}
			return -1
		}, val)
	}
	return val
}
