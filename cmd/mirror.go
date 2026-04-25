package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var srcName, srcEnv, dstName, dstEnv string
	var keys []string
	var overwrite bool
	var prefix string

	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "Mirror keys from one env set into another",
		Long: `Copy key-value pairs from a source env set into a destination env set.
Locked and protected keys in the destination are silently skipped.
Use --prefix to namespace mirrored keys and --overwrite to replace existing values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}

			src, err := store.Load(srcName, srcEnv)
			if err != nil {
				return fmt.Errorf("source not found: %w", err)
			}
			dst, err := store.Load(dstName, dstEnv)
			if err != nil {
				return fmt.Errorf("destination not found: %w", err)
			}

			opts := envset.DefaultMirrorOptions()
			opts.Keys = keys
			opts.Overwrite = overwrite
			opts.Prefix = prefix

			n, err := envset.Mirror(src, dst, opts)
			if err != nil {
				return err
			}

			if err := store.Save(dst); err != nil {
				return err
			}

			mirrored := "all"
			if len(keys) > 0 {
				mirrored = strings.Join(keys, ", ")
			}
			fmt.Printf("Mirrored %d key(s) [%s] from %s/%s → %s/%s\n",
				n, mirrored, srcName, srcEnv, dstName, dstEnv)
			return nil
		},
	}

	cmd.Flags().StringVar(&srcName, "src", "", "Source env set name (required)")
	cmd.Flags().StringVar(&srcEnv, "src-env", "local", "Source environment")
	cmd.Flags().StringVar(&dstName, "dst", "", "Destination env set name (required)")
	cmd.Flags().StringVar(&dstEnv, "dst-env", "local", "Destination environment")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Specific keys to mirror (default: all)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in destination")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to prepend to mirrored key names")
	_ = cmd.MarkFlagRequired("src")
	_ = cmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(cmd)
}
