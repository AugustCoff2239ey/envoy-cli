package envset

import (
	"bufio"
	"fmt"
	"strings"
)

// Import parses a dotenv-formatted string and populates the given EnvSet.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are skipped. Each valid line must be in KEY=VALUE format.
func Import(es *EnvSet, dotenv string) error {
	scanner := bufio.NewScanner(strings.NewReader(dotenv))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip leading "export " if present (shell format compatibility)
		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("line %d: invalid format %q (expected KEY=VALUE)", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

		if err := es.Set(key, value); err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
	}
	return scanner.Err()
}

// ImportFrom parses dotenv content and returns a new EnvSet with the given name and environment.
func ImportFrom(name, environment, dotenv string) (*EnvSet, error) {
	es, err := New(name, environment)
	if err != nil {
		return nil, err
	}
	if err := Import(es, dotenv); err != nil {
		return nil, err
	}
	return es, nil
}
