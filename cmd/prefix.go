package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, prefix string
	var selectedKeys []string
	var strip bool

	cmd := &cobra.Command{
		Use:   "prefix",
		Short: "Add or strip a key prefix in an envset",
		RunE: func(cmd *cobra.Command, args []string) error {
			if prefix == "" {
				return fmt.Errorf("--prefix is required")
			}
			st, err := envset.NewStore(storePath)
			if err != nil {
				return err
			}
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found: %w", name, env, err)
			}

			var n int
			if strip {
				n, err = envset.StripPrefix(es, prefix, selectedKeys)
			} else {
				n, err = envset.AddPrefix(es, prefix, selectedKeys)
			}
			if err != nil {
				return err
			}

			if err := st.Save(es); err != nil {
				return fmt.Errorf("failed to save: %w", err)
			}

			action := "added"
			if strip {
				action = "stripped"
			}
			fmt.Printf("prefix %q %s from %d key(s) in %q (%s)\n",
				prefix, action, n, name, env)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix string (required)")
	cmd.Flags().StringSliceVarP(&selectedKeys, "keys", "k", nil, "Comma-separated keys to target (default: all)")
	cmd.Flags().BoolVarP(&strip, "strip", "s", false, "Strip prefix instead of adding it")
	_ = cmd.MarkFlagRequired("name")

	// normalise any accidental lowercase in prefix flag before use
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		prefix = strings.ToUpper(prefix)
		return nil
	}

	rootCmd.AddCommand(cmd)
}
