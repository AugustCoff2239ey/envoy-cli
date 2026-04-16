package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, transform string
	var keys []string
	var skipLocked bool

	var transformCmd = &cobra.Command{
		Use:   "transform",
		Short: "Apply a transform to env var values",
		Long: `Apply a built-in or no-op transform to env var values in an envset.

Built-in transforms: upper, lower, trim, reverse`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTransform(transform); err != nil {
				return err
			}

			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}

			opts := envset.TransformOptions{
				Keys:       keys,
				SkipLocked: skipLocked,
			}

			if err := envset.Transform(es, transform, nil, opts); err != nil {
				return err
			}

			if err := store.Save(es); err != nil {
				return fmt.Errorf("failed to save envset: %w", err)
			}

			appliedTo := "all keys"
			if len(keys) > 0 {
				appliedTo = strings.Join(keys, ", ")
			}
			fmt.Fprintf(os.Stdout, "Applied %q transform to %s in %q (%s)\n", transform, appliedTo, name, env)
			return nil
		},
	}

	transformCmd.Flags().StringVarP(&name, "name", "n", "", "Envset name (required)")
	transformCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	transformCmd.Flags().StringVarP(&transform, "transform", "t", "", "Transform to apply (upper|lower|trim|reverse)")
	transformCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated keys to transform (default: all)")
	transformCmd.Flags().BoolVar(&skipLocked, "skip-locked", false, "Skip locked keys instead of erroring")
	_ = transformCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(transformCmd)
}

// validTransforms lists the supported built-in transform names.
var validTransforms = map[string]bool{
	"upper":   true,
	"lower":   true,
	"trim":    true,
	"reverse": true,
}

// validateTransform returns an error if the given transform name is not supported.
func validateTransform(transform string) error {
	if transform == "" {
		return fmt.Errorf("--transform is required (upper, lower, trim, reverse)")
	}
	if !validTransforms[transform] {
		return fmt.Errorf("unknown transform %q: must be one of upper, lower, trim, reverse", transform)
	}
	return nil
}
