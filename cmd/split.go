package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var keys string
	var prefix string
	var matchedName string
	var unmatchedName string

	cmd := &cobra.Command{
		Use:   "split <name> <environment>",
		Short: "Split an EnvSet into matched and unmatched subsets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, env := args[0], args[1]
			store := envset.NewStore(storeDir())

			src, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("split: %w", err)
			}

			var opts envset.SplitOptions
			opts.MatchedName = matchedName
			opts.UnmatchedName = unmatchedName

			if keys != "" {
				opts.Keys = strings.Split(keys, ",")
			} else if prefix != "" {
				opts.Predicate = func(k, _ string) bool {
					return strings.HasPrefix(k, prefix)
				}
			} else {
				return fmt.Errorf("split: provide --keys or --prefix")
			}

			res, err := envset.Split(src, opts)
			if err != nil {
				return err
			}

			if err := store.Save(res.Matched); err != nil {
				return fmt.Errorf("split: saving matched: %w", err)
			}
			if err := store.Save(res.Unmatched); err != nil {
				return fmt.Errorf("split: saving unmatched: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "matched   → %s (%d keys)\n", res.Matched.Name, len(res.Matched.Vars))
			fmt.Fprintf(cmd.OutOrStdout(), "unmatched → %s (%d keys)\n", res.Unmatched.Name, len(res.Unmatched.Vars))
			return nil
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to match")
	cmd.Flags().StringVar(&prefix, "prefix", "", "key prefix to match")
	cmd.Flags().StringVar(&matchedName, "matched-name", "", "name for the matched EnvSet")
	cmd.Flags().StringVar(&unmatchedName, "unmatched-name", "", "name for the unmatched EnvSet")

	rootCmd.AddCommand(cmd)
}
