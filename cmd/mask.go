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
	var (
		envName     string
		reveal      int
		keys        []string
		sensitive   bool
		outputJSON  bool
	)

	maskCmd := &cobra.Command{
		Use:   "mask <name>",
		Short: "Display an env set with sensitive values masked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], envName)
			if err != nil {
				return fmt.Errorf("env set not found: %w", err)
			}

			var masked map[string]string
			switch {
			case len(keys) > 0:
				masked = envset.MaskKeys(es, keys, reveal)
			case sensitive:
				masked = envset.MaskSensitive(es, reveal)
			default:
				masked = envset.MaskSensitive(es, reveal)
			}

			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(masked)
			}

			sortedKeys := make([]string, 0, len(masked))
			for k := range masked {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Strings(sortedKeys)
			for _, k := range sortedKeys {
				fmt.Printf("%s=%s\n", k, masked[k])
			}
			return nil
		},
	}

	maskCmd.Flags().StringVarP(&envName, "env", "e", "local", "environment name")
	maskCmd.Flags().IntVarP(&reveal, "reveal", "r", 4, "number of suffix characters to reveal")
	maskCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "specific keys to mask")
	maskCmd.Flags().BoolVarP(&sensitive, "sensitive", "s", false, "mask only sensitive keys (default behaviour)")
	maskCmd.Flags().BoolVar(&outputJSON, "json", false, "output as JSON")

	rootCmd.AddCommand(maskCmd)
}
