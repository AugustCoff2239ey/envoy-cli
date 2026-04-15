package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var parentEnv string
	var childEnv string
	var keys []string

	cmd := &cobra.Command{
		Use:   "inherit <parent-name> <child-name>",
		Short: "Inherit keys from a parent env set into a child, skipping existing keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			parentName := args[0]
			childName := args[1]

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			parent, err := store.Load(parentName, parentEnv)
			if err != nil {
				return fmt.Errorf("parent not found (%s/%s): %w", parentName, parentEnv, err)
			}

			child, err := store.Load(childName, childEnv)
			if err != nil {
				return fmt.Errorf("child not found (%s/%s): %w", childName, childEnv, err)
			}

			result, err := envset.Inherit(parent, child, keys)
			if err != nil {
				return fmt.Errorf("inherit failed: %w", err)
			}

			if err := store.Save(child); err != nil {
				return fmt.Errorf("failed to save child: %w", err)
			}

			if len(result.Inherited) > 0 {
				fmt.Printf("Inherited: %s\n", strings.Join(result.Inherited, ", "))
			} else {
				fmt.Println("No keys inherited.")
			}
			if len(result.Skipped) > 0 {
				fmt.Printf("Skipped (already set): %s\n", strings.Join(result.Skipped, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&parentEnv, "parent-env", "production", "Environment of the parent env set")
	cmd.Flags().StringVar(&childEnv, "child-env", "staging", "Environment of the child env set")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Specific keys to inherit (default: all)")

	rootCmd.AddCommand(cmd)
}
