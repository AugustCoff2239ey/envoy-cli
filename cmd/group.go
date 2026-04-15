package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage key groups within an envset",
	}

	var name, env, format string

	createCmd := &cobra.Command{
		Use:   "create <set-name> <group-name> <key1,key2,...>",
		Short: "Create a named group of keys",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeDir())
			name, env = args[0], resolveEnv(env)
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("load envset: %w", err)
			}
			keys := strings.Split(args[2], ",")
			if err := envset.CreateGroup(es, args[1], keys); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save envset: %w", err)
			}
			fmt.Printf("Group %q created with %d key(s).\n", args[1], len(keys))
			return nil
		},
	}
	createCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment name")

	listCmd := &cobra.Command{
		Use:   "list <set-name>",
		Short: "List all groups in an envset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeDir())
			env = resolveEnv(env)
			es, err := store.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("load envset: %w", err)
			}
			groups := envset.ListGroups(es)
			if format == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(groups)
			}
			if len(groups) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No groups defined.")
				return nil
			}
			for _, g := range groups {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", g.Name, strings.Join(g.Keys, ", "))
			}
			return nil
		},
	}
	listCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment name")
	listCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text|json")

	deleteCmd := &cobra.Command{
		Use:   "delete <set-name> <group-name>",
		Short: "Delete a group from an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeDir())
			env = resolveEnv(env)
			es, err := store.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("load envset: %w", err)
			}
			if err := envset.DeleteGroup(es, args[1]); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save envset: %w", err)
			}
			fmt.Printf("Group %q deleted.\n", args[1])
			return nil
		},
	}
	deleteCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment name")

	groupCmd.AddCommand(createCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(groupCmd)
}
