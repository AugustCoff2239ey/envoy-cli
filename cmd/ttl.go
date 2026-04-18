package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env string
	var ttlStr string

	ttlCmd := &cobra.Command{
		Use:   "ttl",
		Short: "Manage key TTLs in an envset",
	}

	setCmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Set a TTL on a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			ttl, err := time.ParseDuration(ttlStr)
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", ttlStr, err)
			}
			if err := envset.SetTTL(es, args[0], ttl); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("TTL of %s set to %s\n", args[0], ttlStr)
			return nil
		},
	}
	setCmd.Flags().StringVar(&ttlStr, "duration", "", "TTL duration (e.g. 24h, 30m)")
	_ = setCmd.MarkFlagRequired("duration")

	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get remaining TTL for a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			remaining, err := envset.GetTTL(es, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Remaining TTL for %s: %s\n", args[0], remaining.Round(time.Second))
			return nil
		},
	}

	for _, sub := range []*cobra.Command{setCmd, getCmd} {
		sub.Flags().StringVar(&name, "name", "", "EnvSet name")
		sub.Flags().StringVar(&env, "env", "local", "Environment")
		_ = sub.MarkFlagRequired("name")
		ttlCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(ttlCmd)
}
