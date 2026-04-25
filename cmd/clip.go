package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var maxKeys int
	var keys string
	var skipLocked bool

	cmd := &cobra.Command{
		Use:   "clip <name> <env>",
		Short: "Trim an env set to a maximum number of keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, env := args[0], args[1]

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return err
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("clip: %w", err)
			}

			opts := envset.DefaultClipOptions()
			opts.MaxKeys = maxKeys
			opts.SkipLocked = skipLocked
			if keys != "" {
				opts.Keys = strings.Split(keys, ",")
			}

			removed, err := envset.Clip(es, opts)
			if err != nil {
				return err
			}

			if err := store.Save(es); err != nil {
				return fmt.Errorf("clip: save failed: %w", err)
			}

			if len(removed) == 0 {
				fmt.Println("clip: no keys removed")
			} else {
				fmt.Printf("clip: removed %d key(s): %s\n", len(removed), strings.Join(removed, ", "))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxKeys, "max", "m", 10, "Maximum number of keys to retain")
	cmd.Flags().StringVarP(&keys, "keys", "k", "", "Comma-separated list of keys to prioritise keeping")
	cmd.Flags().BoolVar(&skipLocked, "skip-locked", true, "Do not remove locked keys even when over the limit")

	rootCmd.AddCommand(cmd)
}
