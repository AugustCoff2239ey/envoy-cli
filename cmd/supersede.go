package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var skipLocked bool
	var skipProtected bool
	var keys string

	cmd := &cobra.Command{
		Use:   "supersede <source-name> <source-env> <target-name> <target-env>",
		Short: "Overwrite target env vars with values from source",
		Long: `Supersede copies values from a source EnvSet into a target EnvSet,
always overwriting existing values (unlike copy --no-overwrite).
Locked and protected keys are skipped by default.`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName, srcEnv := args[0], args[1]
			dstName, dstEnv := args[2], args[3]

			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}

			src, err := store.Load(srcName, srcEnv)
			if err != nil {
				return fmt.Errorf("source not found: %w", err)
			}
			dst, err := store.Load(dstName, dstEnv)
			if err != nil {
				return fmt.Errorf("target not found: %w", err)
			}

			opts := envset.SupersedeOptions{
				SkipLocked:    skipLocked,
				SkipProtected: skipProtected,
			}
			if keys != "" {
				opts.Keys = strings.Split(keys, ",")
			}

			n, err := envset.Supersede(src, dst, opts)
			if err != nil {
				return err
			}

			if err := store.Save(dst); err != nil {
				return fmt.Errorf("failed to save target: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "superseded %d key(s) into %s/%s\n", n, dstName, dstEnv)
			return nil
		},
	}

	cmd.Flags().BoolVar(&skipLocked, "skip-locked", true, "skip locked keys in target")
	cmd.Flags().BoolVar(&skipProtected, "skip-protected", true, "skip protected keys in target")
	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to supersede (default: all)")

	rootCmd.AddCommand(cmd)
}
