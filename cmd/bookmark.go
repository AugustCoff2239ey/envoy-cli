package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	bookmarkCmd := &cobra.Command{
		Use:   "bookmark",
		Short: "Manage bookmarks on an envset",
	}

	var name, env, format string

	addCmd := &cobra.Command{
		Use:   "add <set> <bookmark>",
		Short: "Add a bookmark to an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			e, err := st.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if err := envset.AddBookmark(e, args[1]); err != nil {
				return err
			}
			return st.Save(e)
		},
	}
	addCmd.Flags().StringVarP(&env, "env", "e", "local", "environment")

	removeCmd := &cobra.Command{
		Use:   "remove <set> <bookmark>",
		Short: "Remove a bookmark from an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			e, err := st.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if err := envset.RemoveBookmark(e, args[1]); err != nil {
				return err
			}
			return st.Save(e)
		},
	}
	removeCmd.Flags().StringVarP(&env, "env", "e", "local", "environment")

	listCmd := &cobra.Command{
		Use:   "list <set>",
		Short: "List bookmarks on an envset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := defaultStore()
			e, err := st.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			bms := envset.ListBookmarks(e)
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(bms)
			}
			if len(bms) == 0 {
				fmt.Println("no bookmarks")
				return nil
			}
			for _, b := range bms {
				fmt.Printf("%-20s -> %s (%s)\n", b.Name, b.SetName, b.Env)
			}
			return nil
		},
	}
	listCmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	listCmd.Flags().StringVar(&format, "format", "text", "output format: text|json")

	bookmarkCmd.AddCommand(addCmd, removeCmd, listCmd)
	_ = name
	rootCmd.AddCommand(bookmarkCmd)
}
