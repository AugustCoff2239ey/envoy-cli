package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var prefix, suffix, envs, exclude string

	cmd := &cobra.Command{
		Use:   "filter <name> <environment>",
		Short: "Filter keys in an env set by prefix, suffix, or exclusion list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, environment := args[0], args[1]

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return err
			}

			es, err := store.Load(name, environment)
			if err != nil {
				return fmt.Errorf("filter: %w", err)
			}

			opts := envset.FilterOptions{
				Prefix: prefix,
				Suffix: suffix,
			}
			if envs != "" {
				opts.Envs = strings.Split(envs, ",")
			}
			if exclude != "" {
				opts.Exclude = strings.Split(exclude, ",")
			}

			out, err := envset.Filter(es, opts)
			if err != nil {
				return err
			}

			if len(out.Vars) == 0 {
				fmt.Println("(no keys matched filter)")
				return nil
			}

			for k, v := range out.Vars {
				fmt.Printf("%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Only include keys with this prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "Only include keys with this suffix")
	cmd.Flags().StringVar(&envs, "envs", "", "Comma-separated environments to match")
	cmd.Flags().StringVar(&exclude, "exclude", "", "Comma-separated keys to exclude")

	rootCmd.AddCommand(cmd)
}
