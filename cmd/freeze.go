package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var environment string

	freezeCmd := &cobra.Command{
		Use:   "freeze <name>",
		Short: "Freeze an envset to prevent modifications",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], environment)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.Freeze(es); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("envset %q (%s) is now frozen\n", es.Name, es.Environment)
			return nil
		},
	}

	unfreezeCmd := &cobra.Command{
		Use:   "unfreeze <name>",
		Short: "Unfreeze an envset to allow modifications",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], environment)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.Unfreeze(es); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("envset %q (%s) is now unfrozen\n", es.Name, es.Environment)
			return nil
		},
	}

	for _, c := range []*cobra.Command{freezeCmd, unfreezeCmd} {
		c.Flags().StringVarP(&environment, "env", "e", "local", "Environment name")
		rootCmd.AddCommand(c)
	}
}
