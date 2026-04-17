package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env string

	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage key aliases in an envset",
	}

	addCmd := &cobra.Command{
		Use:   "add <key> <alias>",
		Short: "Add an alias for a key",
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
			if err := envset.AddAlias(es, args[0], args[1]); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("alias %q -> %q added\n", args[1], args[0])
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <key> <alias>",
		Short: "Remove an alias from a key",
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
			if err := envset.RemoveAlias(es, args[0], args[1]); err != nil {
				return err
			}
			return store.Save(es)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <key>",
		Short: "List all aliases for a key",
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
			aliases := envset.ListAliases(es, args[0])
			if jsonOut, _ := cmd.Flags().GetBool("json"); jsonOut {
				return json.NewEncoder(os.Stdout).Encode(aliases)
			}
			for _, a := range aliases {
				fmt.Println(a)
			}
			return nil
		},
	}
	listCmd.Flags().Bool("json", false, "Output as JSON")

	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd} {
		sub.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
		sub.Flags().StringVarP(&env, "env", "e", "local", "Environment")
		_ = sub.MarkFlagRequired("name")
		aliasCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(aliasCmd)
}
