package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format string
	var purge bool

	cmd := &cobra.Command{
		Use:   "expire <key> [duration]",
		Short: "Set or inspect expiry on an env var key",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeFile())
			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}

			if purge {
				purged, err := envset.PurgeExpired(e)
				if err != nil {
					return err
				}
				if format == "json" {
					json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"purged": purged})
				} else {
					if len(purged) == 0 {
						fmt.Println("No expired keys found.")
					} else {
						for _, k := range purged {
							fmt.Printf("purged: %s\n", k)
						}
					}
				}
				return store.Save(e)
			}

			if len(args) == 0 {
				return fmt.Errorf("key argument required")
			}
			key := args[0]

			if len(args) == 2 {
				d, err := time.ParseDuration(args[1])
				if err != nil {
					return fmt.Errorf("invalid duration %q: %w", args[1], err)
				}
				if err := envset.SetExpiry(e, key, time.Now().Add(d)); err != nil {
					return err
				}
				fmt.Printf("expiry set on %q for %s\n", key, d)
				return store.Save(e)
			}

			t, set, err := envset.GetExpiry(e, key)
			if err != nil {
				return err
			}
			if !set {
				fmt.Printf("%s: no expiry set\n", key)
				return nil
			}
			expired, _ := envset.IsExpired(e, key)
			if format == "json" {
				json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"key": key, "expires_at": t, "expired": expired})
			} else {
				status := "active"
				if expired {
					status = "EXPIRED"
				}
				fmt.Printf("%s: expires %s [%s]\n", key, t.Format(time.RFC3339), status)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "EnvSet name")
	cmd.Flags().StringVar(&env, "env", "local", "Environment")
	cmd.Flags().StringVar(&format, "format", "text", "Output format (text|json)")
	cmd.Flags().BoolVar(&purge, "purge", false, "Remove all expired keys")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
