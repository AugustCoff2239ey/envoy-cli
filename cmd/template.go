package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, tmplFile string
	var inline string

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Render a template using values from an envset",
		Long: `Replace {{KEY}} placeholders in a template string or file
with values from the specified envset.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}

			var tmpl string
			switch {
			case tmplFile != "":
				b, err := os.ReadFile(tmplFile)
				if err != nil {
					return fmt.Errorf("reading template file: %w", err)
				}
				tmpl = string(b)
			case inline != "":
				tmpl = inline
			default:
				return fmt.Errorf("provide --template or --file")
			}

			res, err := envset.RenderTemplate(e, tmpl)
			if err != nil {
				return err
			}

			fmt.Print(res.Rendered)

			if len(res.Unresolved) > 0 {
				fmt.Fprintf(os.Stderr, "\nwarning: unresolved placeholders: %s\n",
					strings.Join(res.Unresolved, ", "))
				sugs := envset.SuggestPlaceholders(e, res.Unresolved)
				for k, s := range sugs {
					if len(s) > 0 {
						fmt.Fprintf(os.Stderr, "  did you mean %q instead of %q?\n", s[0], k)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "envset name (required)")
	cmd.Flags().StringVar(&env, "env", "local", "environment")
	cmd.Flags().StringVar(&tmplFile, "file", "", "path to template file")
	cmd.Flags().StringVar(&inline, "template", "", "inline template string")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
