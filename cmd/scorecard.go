package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "scorecard",
		Short: "Score the quality of an env set",
		Long:  "Evaluate an env set across naming, documentation, hygiene, and security dimensions and return a 0-100 score.",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("load %q (%s): %w", name, env, err)
			}
			res, err := envset.Scorecard(es)
			if err != nil {
				return err
			}
			if jsonOut {
				return printScorecardJSON(res)
			}
			printScorecardText(res)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "env set name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}

func printScorecardText(res *envset.ScorecardResult) {
	fmt.Printf("Score: %d/100  Grade: %s\n", res.Score, res.Grade)
	fmt.Println("Breakdown:")

	categories := make([]string, 0, len(res.Breakdown))
	for k := range res.Breakdown {
		categories = append(categories, k)
	}
	sort.Strings(categories)
	for _, cat := range categories {
		fmt.Printf("  %-16s %d/25\n", cat, res.Breakdown[cat])
	}

	if len(res.Suggestions) > 0 {
		fmt.Println("Suggestions:")
		for _, s := range res.Suggestions {
			fmt.Printf("  • %s\n", s)
		}
	}
}

func printScorecardJSON(res *envset.ScorecardResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(res)
}
