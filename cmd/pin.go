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

	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin or unpin keys in an env set",
	}

	addCmd := &cobra.Command{
		Use:   "add <key>",
		Short: "Pin a key to its current value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("pin add: %w", err)
			}
			user := os.Getenv("USER")
			if user == "" {
				user = "unknown"
			}
			entry, err := envset.PinKey(es, args[0], user)
			if err != nil {
				return fmt.Errorf("pin add: %w", err)
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("pin add: save: %w", err)
			}
			fmt.Printf("pinned %s=%s (by %s)\n", entry.Key, entry.Value, entry.PinnedBy)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <key>",
		Short: "Unpin a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("pin remove: %w", err)
			}
			if err := envset.UnpinKey(es, args[0]); err != nil {
				return fmt.Errorf("pin remove: %w", err)
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("pin remove: save: %w", err)
			}
			fmt.Printf("unpinned %s\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("pin list: %w", err)
			}
			pinned := envset.PinnedKeys(es)
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(pinned)
			}
			if len(pinned) == 0 {
				fmt.Println("no pinned keys")
				return nil
			}
			for _, k := range pinned {
				fmt.Println(k)
			}
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd} {
		sub.Flags().StringVarP(&name, "name", "n", "", "env set name (required)")
		sub.Flags().StringVarP(&env, "env", "e", "local", "environment")
		_ = sub.MarkFlagRequired("name")
	}
	listCmd.Flags().StringVar(&format, "format", "text", "output format: text|json")

	pinCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}
