package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format string

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint an envset for common issues",
		Long:  "Checks an envset for empty values, lowercase keys, oversized values, and shell built-in shadowing.",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("load envset %q (%s): %w", name, env, err)
			}

			result, err := envset.Lint(es)
			if err != nil {
				return fmt.Errorf("lint: %w", err)
			}

			switch format {
			case "json":
				type jsonFinding struct {
					Key      string `json:"key"`
					Message  string `json:"message"`
					Severity string `json:"severity"`
				}
				var out []jsonFinding
				for _, f := range result.Findings {
					out = append(out, jsonFinding{
						Key:      f.Key,
						Message:  f.Message,
						Severity: string(f.Severity),
					})
				}
				if out == nil {
					out = []jsonFinding{}
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(out)
			default:
				fmt.Println(result.Summary())
			}

			if result.HasErrors() {
				os.Exit(1)
			}
			return nil
		},
	}

	lintCmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	lintCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment (local|staging|production)")
	lintCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	_ = lintCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(lintCmd)
}
