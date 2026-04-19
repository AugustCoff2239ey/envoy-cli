package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, pattern, destName string
	var keys []string
	var remove bool

	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract matching keys from an env set into a new set",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			src, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("extract: %w", err)
			}
			if pattern == "" && len(keys) == 0 {
				return fmt.Errorf("extract: provide --pattern or --keys")
			}
			opts := envset.ExtractOptions{
				Pattern:          pattern,
				Keys:             keys,
				RemoveFromSource: remove,
			}
			res, err := envset.Extract(src, opts)
			if err != nil {
				return err
			}
			if destName != "" {
				res.Extracted.Name = destName
			}
			if err := store.Save(res.Extracted); err != nil {
				return fmt.Errorf("extract: save destination: %w", err)
			}
			if remove {
				if err := store.Save(src); err != nil {
					return fmt.Errorf("extract: update source: %w", err)
				}
			}
			fmt.Printf("Extracted %d key(s): %s\n", len(res.Keys), strings.Join(res.Keys, ", "))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Source env set name (required)")
	cmd.Flags().StringVar(&env, "env", "local", "Source environment")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Regex pattern to match keys")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Explicit keys to extract")
	cmd.Flags().StringVar(&destName, "dest", "", "Destination set name (defaults to source name)")
	cmd.Flags().BoolVar(&remove, "remove", false, "Remove extracted keys from source")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
