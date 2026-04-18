package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	dependCmd := &cobra.Command{
		Use:   "depend",
		Short: "Manage key dependencies within an envset",
	}

	var name, env string

	addDepCmd := &cobra.Command{
		Use:   "add <key> <dep>",
		Short: "Add a dependency from key to dep",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.AddDependency(es, args[0], args[1]); err != nil {
				return err
			}
			return st.Save(es)
		},
	}

	removeDepCmd := &cobra.Command{
		Use:   "remove <key> <dep>",
		Short: "Remove a dependency from key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.RemoveDependency(es, args[0], args[1]); err != nil {
				return err
			}
			return st.Save(es)
		},
	}

	listDepCmd := &cobra.Command{
		Use:   "list <key>",
		Short: "List dependencies of a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			deps := envset.GetDependencies(es, args[0])
			if len(deps) == 0 {
				fmt.Println("no dependencies")
				return nil
			}
			fmt.Println(strings.Join(deps, "\n"))
			return nil
		},
	}

	checkDepCmd := &cobra.Command{
		Use:   "check",
		Short: "Check for missing dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			missing := envset.CheckDependencies(es)
			if len(missing) == 0 {
				fmt.Println("all dependencies satisfied")
				return nil
			}
			fmt.Println("missing dependencies:")
			for _, m := range missing {
				fmt.Println(" ", m)
			}
			return fmt.Errorf("unresolved dependencies")
		},
	}

	for _, sub := range []*cobra.Command{addDepCmd, removeDepCmd, listDepCmd, checkDepCmd} {
		sub.Flags().StringVarP(&name, "name", "n", "", "envset name")
		sub.Flags().StringVarP(&env, "env", "e", "local", "environment")
		_ = sub.MarkFlagRequired("name")
		dependCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(dependCmd)
}
