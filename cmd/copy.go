package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var overwrite bool
	var keys string
	var srcEnv, dstEnv string

	copyCmd := &cobra.Command{
		Use:   "copy <src-name> <dst-name>",
		Short: "Copy environment variables from one EnvSet to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName := args[0]
			dstName := args[1]

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			src, err := store.Load(srcName, srcEnv)
			if err != nil {
				return fmt.Errorf("source not found (%s/%s): %w", srcName, srcEnv, err)
			}

			dst, err := store.Load(dstName, dstEnv)
			if err != nil {
				return fmt.Errorf("destination not found (%s/%s): %w", dstName, dstEnv, err)
			}

			var selectedKeys []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						selectedKeys = append(selectedKeys, k)
					}
				}
			}

			opts := envset.CopyOptions{
				Overwrite: overwrite,
				Keys:      selectedKeys,
			}

			n, err := envset.Copy(src, dst, opts)
			if err != nil {
				return fmt.Errorf("copy failed: %w", err)
			}

			if err := store.Save(dst); err != nil {
				return fmt.Errorf("failed to save destination: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Copied %d key(s) from %s/%s to %s/%s\n",
				n, srcName, srcEnv, dstName, dstEnv)
			return nil
		},
	}

	copyCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	copyCmd.Flags().StringVarP(&keys, "keys", "k", "", "Comma-separated list of keys to copy (default: all)")
	copyCmd.Flags().StringVar(&srcEnv, "src-env", "local", "Environment of the source EnvSet")
	copyCmd.Flags().StringVar(&dstEnv, "dst-env", "local", "Environment of the destination EnvSet")

	rootCmd.AddCommand(copyCmd)
}
