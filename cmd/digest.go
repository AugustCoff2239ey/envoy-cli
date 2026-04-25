package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format, keys string

	cmd := &cobra.Command{
		Use:   "digest",
		Short: "Compute a deterministic SHA-256 digest of an env set",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("digest: %w", err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					k = strings.TrimSpace(k)
					if k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			res, err := envset.Digest(es, keyList)
			if err != nil {
				return err
			}

			if format == "json" {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			}

			fmt.Printf("digest: %s\n", res.Digest)
			for k, v := range res.Entries {
				fmt.Printf("  %s: %s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "env set name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "environment")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text|json")
	cmd.Flags().StringVarP(&keys, "keys", "k", "", "comma-separated list of keys to include")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
