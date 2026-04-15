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
	var name, env, format string

	cmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Expand ${KEY} and ${KEY:default} references within an env set",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}

			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("interpolate: %w", err)
			}

			res, err := envset.Interpolate(e)
			if err != nil {
				return err
			}

			switch strings.ToLower(format) {
			case "json":
				return printInterpolateJSON(res)
			default:
				return printInterpolateText(res)
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "Environment (local/staging/production)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}

func printInterpolateText(res envset.InterpolateResult) error {
	for k, v := range res.Resolved {
		fmt.Printf("%s=%s\n", k, v)
	}
	if len(res.Unresolved) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "unresolved references: %s\n",
			strings.Join(res.Unresolved, ", "))
	}
	return nil
}

func printInterpolateJSON(res envset.InterpolateResult) error {
	out := map[string]interface{}{
		"resolved":   res.Resolved,
		"unresolved": res.Unresolved,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
