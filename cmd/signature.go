package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "Sign and verify EnvSet integrity using HMAC-SHA256",
	}

	var passphrase string

	signCmd := &cobra.Command{
		Use:   "sign <name> <environment>",
		Short: "Sign an EnvSet with a passphrase",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if err := envset.Sign(es, passphrase); err != nil {
				return fmt.Errorf("sign: %w", err)
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "signed %s/%s\n", args[0], args[1])
			return nil
		},
	}
	signCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for HMAC signing (required)")
	_ = signCmd.MarkFlagRequired("passphrase")

	verifyCmd := &cobra.Command{
		Use:   "verify <name> <environment>",
		Short: "Verify the signature of an EnvSet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if err := envset.VerifySignature(es, passphrase); err != nil {
				fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "signature OK for %s/%s\n", args[0], args[1])
			return nil
		},
	}
	verifyCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase used when signing (required)")
	_ = verifyCmd.MarkFlagRequired("passphrase")

	clearCmd := &cobra.Command{
		Use:   "clear <name> <environment>",
		Short: "Remove the stored signature from an EnvSet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			envset.ClearSignature(es)
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "signature cleared for %s/%s\n", args[0], args[1])
			return nil
		},
	}

	signatureCmd.AddCommand(signCmd, verifyCmd, clearCmd)
	rootCmd.AddCommand(signatureCmd)
}
