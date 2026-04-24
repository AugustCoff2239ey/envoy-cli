package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "crossref <base-name> <target-name>",
		Short: "Cross-reference two env sets to show shared, exclusive, and mismatched keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}

			base, err := store.Load(args[0])
			if err != nil {
				return fmt.Errorf("base set %q not found: %w", args[0], err)
			}
			target, err := store.Load(args[1])
			if err != nil {
				return fmt.Errorf("target set %q not found: %w", args[1], err)
			}

			result, err := envset.CrossRef(base, target)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printCrossRefJSON(result)
			}
			printCrossRefText(args[0], args[1], result)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(cmd)
}

func printCrossRefText(baseName, targetName string, r *envset.CrossRefResult) {
	fmt.Printf("Cross-reference: %s ↔ %s\n\n", baseName, targetName)
	fmt.Printf("Shared keys (%d):         %s\n", len(r.SharedKeys), strings.Join(r.SharedKeys, ", "))
	fmt.Printf("Only in base (%d):        %s\n", len(r.OnlyInBase), strings.Join(r.OnlyInBase, ", "))
	fmt.Printf("Only in target (%d):      %s\n", len(r.OnlyInTarget), strings.Join(r.OnlyInTarget, ", "))
	fmt.Printf("Value matches (%d):       %s\n", len(r.ValueMatches), strings.Join(r.ValueMatches, ", "))
	fmt.Printf("Value mismatches (%d):    %s\n", len(r.ValueMismatches), strings.Join(r.ValueMismatches, ", "))
}

func printCrossRefJSON(r *envset.CrossRefResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
