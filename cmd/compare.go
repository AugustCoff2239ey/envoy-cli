package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"envoy-cli/internal/envset"
)

func init() {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "compare <name> <env> <other-name> <other-env>",
		Short: "Compare two env sets key by key",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			base, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("base: %w", err)
			}
			target, err := store.Load(args[2], args[3])
			if err != nil {
				return fmt.Errorf("target: %w", err)
			}
			result, err := envset.Compare(base, target)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printCompareJSON(result)
			}
			printCompareText(result)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(cmd)
}

func printCompareText(r *envset.CompareResult) {
	print := func(label string, keys []string) {
		if len(keys) == 0 {
			return
		}
		fmt.Printf("%s:\n", label)
		for _, k := range keys {
			fmt.Printf("  %s\n", k)
		}
	}
	print("matching", r.Matching)
	print("mismatched", r.Mismatched)
	print("only in base", r.OnlyInBase)
	print("only in target", r.OnlyInTarget)
	if len(r.Matching)+len(r.Mismatched)+len(r.OnlyInBase)+len(r.OnlyInTarget) == 0 {
		fmt.Println("sets are identical")
	}
}

func printCompareJSON(r *envset.CompareResult) error {
	return json.NewEncoder(os.Stdout).Encode(r)
}
