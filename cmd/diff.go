package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy-cli/internal/envset"
	"github.com/spf13/cobra"
)

var diffFormat string

var diffCmd = &cobra.Command{
	Use:   "diff <base-env> <target-env>",
	Short: "Show differences between two environment sets",
	Long: `Compare two environment sets and display added, removed, and changed keys.

Example:
  envoy diff local staging
  envoy diff staging production --format json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseEnv := args[0]
		targetEnv := args[1]

		store, err := envset.NewStore(defaultStorePath())
		if err != nil {
			return fmt.Errorf("failed to open store: %w", err)
		}

		base, err := store.Load(baseEnv)
		if err != nil {
			return fmt.Errorf("failed to load base env %q: %w", baseEnv, err)
		}

		target, err := store.Load(targetEnv)
		if err != nil {
			return fmt.Errorf("failed to load target env %q: %w", targetEnv, err)
		}

		result := envset.Diff(base, target)

		if len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
			return nil
		}

		switch diffFormat {
		case "json":
			printDiffJSON(cmd, result)
		default:
			printDiffText(cmd, baseEnv, targetEnv, result)
		}

		return nil
	},
}

func printDiffText(cmd *cobra.Command, base, target string, result envset.DiffResult) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Diff: %s → %s\n\n", base, target)
	for _, k := range result.Added {
		fmt.Fprintf(out, "  + %s\n", k)
	}
	for _, k := range result.Removed {
		fmt.Fprintf(out, "  - %s\n", k)
	}
	for _, c := range result.Changed {
		fmt.Fprintf(out, "  ~ %s: %q → %q\n", c.Key, c.OldValue, c.NewValue)
	}
}

func printDiffJSON(cmd *cobra.Command, result envset.DiffResult) {
	out := cmd.OutOrStdout()
	data, err := result.MarshalJSON()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal diff to JSON:", err)
		return
	}
	fmt.Fprintln(out, string(data))
}

func init() {
	diffCmd.Flags().StringVar(&diffFormat, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(diffCmd)
}
