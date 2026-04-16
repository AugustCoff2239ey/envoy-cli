package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var env string

	protectCmd := &cobra.Command{
		Use:   "protect <name> <key>",
		Short: "Protect a key from modification or deletion",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, key := args[0], args[1]
			s, err := defaultStore()
			if err != nil {
				return err
			}
			e, err := s.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.ProtectKey(e, key); err != nil {
				return err
			}
			if err := s.Save(e); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "key %q protected in %s/%s\n", key, name, env)
			return nil
		},
	}

	unprotectCmd := &cobra.Command{
		Use:   "unprotect <name> <key>",
		Short: "Remove protection from a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, key := args[0], args[1]
			s, err := defaultStore()
			if err != nil {
				return err
			}
			e, err := s.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			if err := envset.UnprotectKey(e, key); err != nil {
				if errors.Is(err, envset.ErrNotProtected) {
					return fmt.Errorf("key %q is not protected", key)
				}
				return err
			}
			if err := s.Save(e); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "key %q unprotected in %s/%s\n", key, name, env)
			return nil
		},
	}

	listProtectedCmd := &cobra.Command{
		Use:   "list-protected <name>",
		Short: "List all protected keys in an envset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := defaultStore()
			if err != nil {
				return err
			}
			e, err := s.Load(args[0], env)
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			keys := envset.ProtectedKeys(e)
			if len(keys) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no protected keys")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(keys, "\n"))
			return nil
		},
	}

	for _, c := range []*cobra.Command{protectCmd, unprotectCmd, listProtectedCmd} {
		c.Flags().StringVarP(&env, "env", "e", "local", "environment name")
		rootCmd.AddCommand(c)
	}
}
