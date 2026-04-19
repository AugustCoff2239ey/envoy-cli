package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var key, value string
	var stopOnError bool

	chainCmd := &cobra.Command{
		Use:   "chain [name1,name2,...] --env ENV",
		Short: "Apply a key=value set operation across multiple envsets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, _ := cmd.Flags().GetString("env")
			names := strings.Split(args[0], ",")

			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}

			var sets []*envset.EnvSet
			for _, name := range names {
				name = strings.TrimSpace(name)
				es, err := store.Load(name, env)
				if err != nil {
					return fmt.Errorf("chain: could not load %q: %w", name, err)
				}
				sets = append(sets, es)
			}

			result, err := envset.Chain(sets, func(es *envset.EnvSet) error {
				if setErr := es.Set(key, value); setErr != nil {
					return setErr
				}
				return store.Save(es)
			}, stopOnError)

			if err != nil {
				return err
			}

			fmt.Printf("Applied: %s\n", strings.Join(result.Applied, ", "))
			if len(result.Skipped) > 0 {
				fmt.Printf("Skipped: %s\n", strings.Join(result.Skipped, ", "))
			}
			return nil
		},
	}

	chainCmd.Flags().String("env", "local", "environment name")
	chainCmd.Flags().StringVar(&key, "key", "", "key to set (required)")
	chainCmd.Flags().StringVar(&value, "value", "", "value to set (required)")
	chainCmd.Flags().BoolVar(&stopOnError, "stop-on-error", false, "halt chain on first error")
	_ = chainCmd.MarkFlagRequired("key")
	_ = chainCmd.MarkFlagRequired("value")

	rootCmd.AddCommand(chainCmd)
}
