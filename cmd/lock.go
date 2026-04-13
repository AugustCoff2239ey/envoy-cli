package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, lockedBy string
	var unlock, list, jsonOut bool

	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Lock or unlock keys in an envset to prevent modification",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeDir())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found: %w", name, env, err)
			}

			if list {
				entries := envset.LockedKeys(es)
				if jsonOut {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "  ")
					return enc.Encode(entries)
				}
				if len(entries) == 0 {
					fmt.Println("No locked keys.")
					return nil
				}
				for _, e := range entries {
					fmt.Printf("  %s  (locked by: %s, at: %s)\n", e.Key, e.LockedBy, e.LockedAt.Format("2006-01-02 15:04:05"))
				}
				return nil
			}

			if len(args) == 0 {
				return fmt.Errorf("at least one key is required")
			}

			for _, key := range args {
				if unlock {
					if err := envset.UnlockKey(es, key); err != nil {
						return fmt.Errorf("unlock %q: %w", key, err)
					}
					fmt.Printf("Unlocked key: %s\n", key)
				} else {
					if err := envset.LockKey(es, key, lockedBy); err != nil {
						return fmt.Errorf("lock %q: %w", key, err)
					}
					fmt.Printf("Locked key: %s\n", key)
				}
			}

			return store.Save(es)
		},
	}

	lockCmd.Flags().StringVarP(&name, "name", "n", "", "envset name (required)")
	lockCmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	lockCmd.Flags().StringVar(&lockedBy, "by", "user", "identifier of who is locking the key")
	lockCmd.Flags().BoolVar(&unlock, "unlock", false, "unlock the specified keys instead")
	lockCmd.Flags().BoolVar(&list, "list", false, "list all locked keys")
	lockCmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON (used with --list)")
	_ = lockCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(lockCmd)
}
