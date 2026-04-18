package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var environment string

	readonlyCmd := &cobra.Command{
		Use:   "readonly <name>",
		Short: "Manage read-only status of an envset",
	}

	markCmd := &cobra.Command{
		Use:   "mark <name>",
		Short: "Mark an envset as read-only",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], environment)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.MarkReadonly(es); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("envset %q (%s) marked as read-only\n", es.Name, es.Environment)
			return nil
		},
	}

	unmarkCmd := &cobra.Command{
		Use:   "unmark <name>",
		Short: "Remove read-only flag from an envset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], environment)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.UnmarkReadonly(es); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("envset %q (%s) is now writable\n", es.Name, es.Environment)
			return nil
		},
	}

	for _, sub := range []*cobra.Command{markCmd, unmarkCmd} {
		sub.Flags().StringVarP(&environment, "env", "e", "local", "Environment name")
		readonlyCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(readonlyCmd)
}
