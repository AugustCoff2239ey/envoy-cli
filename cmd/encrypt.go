package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var environment string
	var passphrase string
	var decrypt bool

	encryptCmd := &cobra.Command{
		Use:   "encrypt <name> <key>",
		Short: "Encrypt or decrypt a single environment variable value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			key := args[1]

			if passphrase == "" {
				passphrase = os.Getenv("ENVOY_PASSPHRASE")
			}
			if passphrase == "" {
				return fmt.Errorf("passphrase is required (use --passphrase or set ENVOY_PASSPHRASE)")
			}

			store := envset.NewStore(storeDir())
			es, err := store.Load(name, environment)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}

			val, ok := es.Vars[key]
			if !ok {
				return fmt.Errorf("key %q not found in %s/%s", key, name, environment)
			}

			var result string
			if decrypt {
				result, err = envset.Decrypt(val, passphrase)
				if err != nil {
					return fmt.Errorf("decrypt: %w", err)
				}
			} else {
				result, err = envset.Encrypt(val, passphrase)
				if err != nil {
					return fmt.Errorf("encrypt: %w", err)
				}
			}

			es.Vars[key] = result
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			action := "Encrypted"
			if decrypt {
				action = "Decrypted"
			}
			fmt.Printf("%s value for key %q in %s/%s\n", action, key, name, environment)
			return nil
		},
	}

	encryptCmd.Flags().StringVarP(&environment, "env", "e", "local", "environment name")
	encryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "encryption passphrase")
	encryptCmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "decrypt the value instead of encrypting")

	rootCmd.AddCommand(encryptCmd)
}
