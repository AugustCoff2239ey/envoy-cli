package envset

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// ObfuscateOptions controls how values are obfuscated.
type ObfuscateOptions struct {
	// Keys to obfuscate; if empty, all keys are obfuscated.
	Keys []string
	// UseHash replaces the value with a SHA-256 prefix instead of asterisks.
	UseHash bool
}

// DefaultObfuscateOptions returns sensible defaults.
func DefaultObfuscateOptions() ObfuscateOptions {
	return ObfuscateOptions{
		UseHash: false,
	}
}

// Obfuscate replaces variable values with obscured representations.
// If opts.Keys is non-empty, only those keys are obfuscated.
func Obfuscate(es *EnvSet, opts ObfuscateOptions) (*EnvSet, error) {
	if es == nil {
		return nil, fmt.Errorf("obfuscate: nil EnvSet")
	}

	target := &EnvSet{
		Name:        es.Name,
		Environment: es.Environment,
		Vars:        make(map[string]string, len(es.Vars)),
		Meta:        es.Meta,
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	for k, v := range es.Vars {
		if len(keySet) == 0 || keySet[k] {
			target.Vars[k] = obfuscateValue(v, opts.UseHash)
		} else {
			target.Vars[k] = v
		}
	}

	return target, nil
}

func obfuscateValue(v string, useHash bool) string {
	if v == "" {
		return ""
	}
	if useHash {
		sum := sha256.Sum256([]byte(v))
		return fmt.Sprintf("sha256:%x", sum[:4])
	}
	return strings.Repeat("*", min8(len(v)))
}

func min8(n int) int {
	if n > 8 {
		return 8
	}
	return n
}
