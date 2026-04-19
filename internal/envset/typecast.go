package envset

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeCastResult holds the result of a typecast operation.
type TypeCastResult struct {
	Key      string
	Original string
	Casted   string
	Type     string
}

// TypeCast converts the values of specified keys to the given type.
// Supported types: "int", "float", "bool", "upper", "lower".
func TypeCast(es *EnvSet, castType string, keys []string) ([]TypeCastResult, error) {
	if es == nil {
		return nil, fmt.Errorf("typecast: envset is nil")
	}
	if castType == "" {
		return nil, fmt.Errorf("typecast: cast type must not be empty")
	}

	targets := keys
	if len(targets) == 0 {
		for k := range es.Vars {
			targets = append(targets, k)
		}
	}

	var results []TypeCastResult
	for _, k := range targets {
		v, ok := es.Vars[k]
		if !ok {
			return nil, fmt.Errorf("typecast: key %q not found", k)
		}
		casted, err := castValue(v, castType)
		if err != nil {
			return nil, fmt.Errorf("typecast: key %q: %w", k, err)
		}
		es.Vars[k] = casted
		results = append(results, TypeCastResult{
			Key:      k,
			Original: v,
			Casted:   casted,
			Type:     castType,
		})
	}
	return results, nil
}

func castValue(v, castType string) (string, error) {
	switch castType {
	case "int":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int", v)
		}
		return strconv.Itoa(int(f)), nil
	case "float":
		_, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float", v)
		}
		return strings.TrimSpace(v), nil
	case "bool":
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to bool", v)
		}
		return strconv.FormatBool(b), nil
	case "upper":
		return strings.ToUpper(v), nil
	case "lower":
		return strings.ToLower(v), nil
	default:
		return "", fmt.Errorf("unsupported cast type %q", castType)
	}
}
