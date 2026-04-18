package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format string

	cmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify keys in an env set by category (secret, database, network, feature, general)",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("classify: %w", err)
			}
			report, err := envset.Classify(es)
			if err != nil {
				return err
			}
			if format == "json" {
				return printClassifyJSON(report)
			}
			return printClassifyText(report)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text|json")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}

func printClassifyText(report *envset.ClassifyReport) error {
	categories := make([]string, 0, len(report.ByCategory))
	for c := range report.ByCategory {
		categories = append(categories, c)
	}
	sort.Strings(categories)
	for _, cat := range categories {
		keys := report.ByCategory[cat]
		sort.Strings(keys)
		fmt.Printf("[%s]\n", cat)
		for _, k := range keys {
			fmt.Printf("  %s\n", k)
		}
	}
	return nil
}

func printClassifyJSON(report *envset.ClassifyReport) error {
	enc := json.NewEncoder(rootCmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(report.ByCategory)
}
