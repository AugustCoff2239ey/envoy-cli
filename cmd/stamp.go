package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/envset"
)

func init() {
	var name, env, format, prefix, suffix string
	var keys []string

	stampCmd := &cobra.Command{
		Use:   "stamp",
		Short: "Write the current UTC timestamp into one or more env-var values",
		Example: `  envoy stamp --name myapp --env production
  envoy stamp --name myapp --env staging --keys BUILT_AT,DEPLOYED_AT
  envoy stamp --name myapp --env local --format "2006-01-02" --prefix "built:"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return fmt.Errorf("stamp: %w", err)
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("stamp: %w", err)
			}

			opts := envset.DefaultStampOptions()
			if format != "" {
				opts.Format = format
			}
			opts.Prefix = prefix
			opts.Suffix = suffix

			if len(keys) > 0 {
				for _, raw := range keys {
					for _, k := range strings.Split(raw, ",") {
						k = strings.TrimSpace(k)
						if k != "" {
							opts.Keys = append(opts.Keys, k)
						}
					}
				}
			}

			stamped, err := envset.Stamp(es, opts)
			if err != nil {
				return fmt.Errorf("stamp: %w", err)
			}

			if err := store.Save(es); err != nil {
				return fmt.Errorf("stamp: save: %w", err)
			}

			if len(stamped) == 0 {
				fmt.Println("no keys were stamped (all may be locked or protected)")
				return nil
			}
			for k, v := range stamped {
				fmt.Printf("  %-30s = %s\n", k, v)
			}
			fmt.Printf("stamped %d key(s)\n", len(stamped))
			return nil
		},
	}

	stampCmd.Flags().StringVar(&name, "name", "", "EnvSet name (required)")
	stampCmd.Flags().StringVar(&env, "env", "local", "Environment (local|staging|production)")
	stampCmd.Flags().StringArrayVar(&keys, "keys", nil, "Comma-separated keys to stamp (default: all)")
	stampCmd.Flags().StringVar(&format, "format", "", "Go time layout string (default: RFC3339)")
	stampCmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to prepend to the timestamp")
	stampCmd.Flags().StringVar(&suffix, "suffix", "", "Suffix to append to the timestamp")
	_ = stampCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(stampCmd)
}
