package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env,	Use:   "typecast",
		Short: "Cast environment variable values to a specified type",
		Run)
			if err != nil {
				return err
			}
			es, err := s.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}
			results, err := envset.TypeCast(es, castType, keys)
			if err != nil {
				return err
			}
			if err := s.Save(es); err != nil {
				return err
			}
			if output == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %q -> %q (%s)\n", r.Key, r.Original, r.Casted, r.Type)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	cmd.Flags().StringVarP(&castType, "type", "t", "", "Cast type: int, float, bool, upper, lower (required)")
	cmd.Flags().StringArrayVarP(&keys, "key", "k", nil, "Keys to cast (default: all)")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text, json")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")

	_ = strings.ToLower // suppress unused import
	rootCmd.AddCommand(cmd)
}
