package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/envset"
)

func init() {
	var overwrite bool
	var keys string

	cmd := &cobra.Command{
		Use:   "promote <name> <source-env> <target-env>",
		Short: "Promote variables from one environment to another",
		Long: `Copies variables from the source environment into the target environment.
By default existing keys in the target are preserved unless --overwrite is set.
Use --keys to restrict promotion to a comma-separated list of variable names.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, srcEnv, dstEnv := args[0], args[1], args[2]

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return err
			}

			source, err := store.Load(name, srcEnv)
			if err != nil {
				return fmt.Errorf("source %s/%s not found: %w", name, srcEnv, err)
			}

			target, err := store.Load(name, dstEnv)
			if err != nil {
				// Target may not exist yet; create an empty one.
				target, err = envset.New(name, dstEnv)
				if err != nil {
					return err
				}
			}

			var selectedKeys []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						selectedKeys = append(selectedKeys, k)
					}
				}
			}

			result, err := envset.Promote(source, target, envset.PromoteOptions{
				Overwrite: overwrite,
				Keys:      selectedKeys,
			})
			if err != nil {
				return err
			}

			if err := store.Save(result); err != nil {
				return fmt.Errorf("saving promoted envset: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "promoted %s from %s → %s\n", name, srcEnv, dstEnv)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in the target environment")
	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to promote (default: all)")

	rootCmd.AddCommand(cmd)
}
