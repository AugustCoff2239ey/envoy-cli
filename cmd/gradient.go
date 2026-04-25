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
	var names []string
	var keys []string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "gradient",
		Short: "Trace how key values change across multiple envsets",
		Example: `  envoy gradient --names app:local,app:staging,app:production
  envoy gradient --names app:local,app:production --keys DB_HOST,LOG_LEVEL
  envoy gradient --names app:local,app:staging --format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(names) == 0 {
				return fmt.Errorf("--names is required (e.g. app:local,app:staging)")
			}
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			var sets []*envset.EnvSet
			for _, entry := range names {
				parts := strings.SplitN(entry, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid name format %q, expected name:environment", entry)
				}
				es, err := store.Load(parts[0], parts[1])
				if err != nil {
					return fmt.Errorf("load %q: %w", entry, err)
				}
				sets = append(sets, es)
			}
			results, err := envset.Gradient(sets, keys)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return json.NewEncoder(os.Stdout).Encode(results)
			}
			for _, r := range results {
				uniform := "varies"
				if r.Uniform {
					uniform = "uniform"
				}
				fmt.Printf("[%s] (%s)\n", r.Key, uniform)
				for _, step := range r.Steps {
					fmt.Printf("  %-16s = %s\n", step.Environment, step.Value)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&names, "names", nil, "Comma-separated list of name:environment pairs")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Comma-separated list of keys to trace (default: all)")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}
