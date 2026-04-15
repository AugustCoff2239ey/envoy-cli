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
	var caseSensitive, regex, keysOnly, valuesOnly bool

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search keys and values in an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("env set %q (%s) not found", name, env)
			}
			opts := envset.SearchOptions{
				CaseSensitive: caseSensitive,
				Regex:         regex,
				KeysOnly:      keysOnly,
				ValuesOnly:    valuesOnly,
			}
			results, err := envset.Search(es, query, opts)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				fmt.Fprintln(os.Stderr, "no matches found")
				return nil
			}
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(results)
			}
			for _, r := range results {
				fmt.Printf("%-30s = %s  [match: %s]\n", r.Key, r.Value, r.Field)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "env set name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text|json")
	cmd.Flags().BoolVar(&caseSensitive, "case-sensitive", false, "enable case-sensitive matching")
	cmd.Flags().BoolVar(&regex, "regex", false, "treat query as a regular expression")
	cmd.Flags().BoolVar(&keysOnly, "keys-only", false, "search keys only")
	cmd.Flags().BoolVar(&valuesOnly, "values-only", false, "search values only")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
