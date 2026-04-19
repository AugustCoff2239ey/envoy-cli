package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name string
	var env string

	cmd := &cobra.Command{
		Use:   "reorder --name NAME --env ENV KEY [KEY ...]",
		Short: "Set the export/display order of keys in an env set",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("reorder: %w", err)
			}
			res, err := envset.Reorder(es, args)
			if err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Printf("Key order updated: %s\n", strings.Join(res.Keys, ", "))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "EnvSet name (required)")
	cmd.Flags().StringVar(&env, "env", "local", "Environment")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
