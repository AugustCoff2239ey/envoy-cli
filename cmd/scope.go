package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage key scopes within an env set",
	}

	var scopeFormat string

	createCmd := &cobra.Command{
		Use:   "create <name> <env-set> <KEY,...>",
		Short: "Create a new scope with the given keys",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopeName, setName, keyList := args[0], args[1], args[2]
			keys := strings.Split(keyList, ",")
			store := envset.NewStore(defaultStorePath())
			es, err := store.Load(setName)
			if err != nil {
				return fmt.Errorf("load %q: %w", setName, err)
			}
			if err := envset.CreateScope(es, scopeName, keys); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Printf("scope %q created in %q\n", scopeName, setName)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env-set>",
		Short: "List all scopes in an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(defaultStorePath())
			es, err := store.Load(args[0])
			if err != nil {
				return fmt.Errorf("load %q: %w", args[0], err)
			}
			names := envset.ListScopes(es)
			sort.Strings(names)
			if scopeFormat == "json" {
				return json.NewEncoder(os.Stdout).Encode(names)
			}
			if len(names) == 0 {
				fmt.Println("no scopes defined")
				return nil
			}
			for _, n := range names {
				keys, _ := envset.GetScope(es, n)
				fmt.Printf("%-20s %s\n", n, strings.Join(keys, ", "))
			}
			return nil
		},
	}
	listCmd.Flags().StringVarP(&scopeFormat, "format", "f", "text", "Output format: text|json")

	deleteCmd := &cobra.Command{
		Use:   "delete <scope> <env-set>",
		Short: "Delete a scope from an env set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopeName, setName := args[0], args[1]
			store := envset.NewStore(defaultStorePath())
			es, err := store.Load(setName)
			if err != nil {
				return fmt.Errorf("load %q: %w", setName, err)
			}
			if err := envset.DeleteScope(es, scopeName); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Printf("scope %q deleted from %q\n", scopeName, setName)
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <scope> <env-set>",
		Short: "Show the keys belonging to a scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopeName, setName := args[0], args[1]
			store := envset.NewStore(defaultStorePath())
			es, err := store.Load(setName)
			if err != nil {
				return fmt.Errorf("load %q: %w", setName, err)
			}
			keys, err := envset.GetScope(es, scopeName)
			if err != nil {
				return fmt.Errorf("scope %q not found in %q", scopeName, setName)
			}
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}

	scopeCmd.AddCommand(createCmd, listCmd, deleteCmd, showCmd)
	rootCmd.AddCommand(scopeCmd)
}
