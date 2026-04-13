package envset

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exported env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Export serializes an EnvSet into the specified format string.
func Export(es *EnvSet, format Format) (string, error) {
	switch format {
	case FormatDotenv:
		return exportDotenv(es), nil
	case FormatExport:
		return exportShell(es), nil
	case FormatJSON:
		return exportJSON(es), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", format)
	}
}

func sortedKeys(es *EnvSet) []string {
	keys := make([]string, 0, len(es.Vars))
	for k := range es.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func exportDotenv(es *EnvSet) string {
	var sb strings.Builder
	for _, k := range sortedKeys(es) {
		fmt.Fprintf(&sb, "%s=%s\n", k, es.Vars[k])
	}
	return sb.String()
}

func exportShell(es *EnvSet) string {
	var sb strings.Builder
	for _, k := range sortedKeys(es) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, es.Vars[k])
	}
	return sb.String()
}

func exportJSON(es *EnvSet) string {
	keys := sortedKeys(es)
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, es.Vars[k], comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
