package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format string

	labelCmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels on an envset",
	}

	addCmd := &cobra.Command{
		Use:   "add <key> <value>",
		Short: "Add a label to an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.AddLabel(es, args[0], args[1]); err != nil {
				return err
			}
			return store.Save(es)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <key>",
		Short: "Remove a label from an envset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.RemoveLabel(es, args[0]); err != nil {
				return err
			}
			return store.Save(es)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List labels on an envset",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			labels := envset.ListLabels(es)
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(labels)
			}
			for k, v := range labels {
				fmt.Printf("%s=%s\n", k, v)
			}
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd} {
		sub.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
		sub.Flags().StringVarP(&env, "env", "e", "local", "Environment")
		_ = sub.MarkFlagRequired("name")
	}
	listCmd.Flags().StringVar(&format, "format", "text", "Output format: text|json")

	labelCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(labelCmd)
}
