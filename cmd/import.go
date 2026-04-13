package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var environment string
	var filePath string
	var format string
	var overwrite bool

	importCmd := &cobra.Command{
		Use:   "import <name>",
		Short: "Import environment variables into an envset from a file",
		Long: `Import environment variables into a named envset from a .env, shell, or JSON file.

Supported formats: dotenv, shell, json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			store, err := envset.NewStore("")
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			var set *envset.EnvSet

			existing, loadErr := store.Load(name, environment)
			if loadErr == nil && !overwrite {
				set = existing
			} else {
				set, err = envset.New(name, environment)
				if err != nil {
					return fmt.Errorf("failed to create envset: %w", err)
				}
			}

			if filePath == "" {
				if err := envset.Import(set, os.Stdin, format); err != nil {
					return fmt.Errorf("failed to import from stdin: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Imported into envset '%s' [%s] from stdin\n", name, environment)
			} else {
				if err := envset.ImportFrom(set, filePath, format); err != nil {
					return fmt.Errorf("failed to import from file: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Imported into envset '%s' [%s] from %s\n", name, environment, filePath)
			}

			if err := store.Save(set); err != nil {
				return fmt.Errorf("failed to save envset: %w", err)
			}

			return nil
		},
	}

	importCmd.Flags().StringVarP(&environment, "env", "e", "local", "Target environment (local, staging, production)")
	importCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to input file (reads from stdin if not specified)")
	importCmd.Flags().StringVar(&format, "format", "dotenv", "Input format: dotenv, shell, json")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing envset instead of merging")

	rootCmd.AddCommand(importCmd)
}
