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
	var name, env, format string
	var quotes, leadOnly, trailOnly bool
	var keys []string

	cmd := &cobra.Command{
		Use:   "trim",
		Short: "Trim whitespace (and optionally quotes) from env var values",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}

			opts := envset.DefaultTrimOptions()
			if cmd.Flags().Changed("no-leading") {
				opts.LeadingSpace = false
			}
			if cmd.Flags().Changed("no-trailing") {
				opts.TrailingSpace = false
			}
			if leadOnly {
				opts.TrailingSpace = false
			}
			if trailOnly {
				opts.LeadingSpace = false
			}
			opts.Quotes = quotes
			opts.Keys = keys

			results, err := envset.Trim(es, opts)
			if err != nil {
				return err
			}

			if err := store.Save(es); err != nil {
				return fmt.Errorf("failed to save envset: %w", err)
			}

			switch strings.ToLower(format) {
			case "json":
				return json.NewEncoder(os.Stdout).Encode(results)
			default:
				changed := 0
				for _, r := range results {
					if r.Changed {
						fmt.Printf("  trimmed %s: %q -> %q\n", r.Key, r.OldVal, r.NewVal)
						changed++
					}
				}
				if changed == 0 {
					fmt.Println("no values required trimming")
				} else {
					fmt.Printf("%d value(s) trimmed and saved.\n", changed)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "envset name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text|json")
	cmd.Flags().BoolVar(&quotes, "quotes", false, "also strip surrounding quotes")
	cmd.Flags().BoolVar(&leadOnly, "leading-only", false, "trim leading whitespace only")
	cmd.Flags().BoolVar(&trailOnly, "trailing-only", false, "trim trailing whitespace only")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated keys to trim (default: all)")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
