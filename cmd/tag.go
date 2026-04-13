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
	var environment string
	var outputJSON bool

	tagCmd := &cobra.Command{
		Use:   "tag <name> <add|remove|list> [tag-name]",
		Short: "Manage tags on an envset",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			action := strings.ToLower(args[1])

			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return fmt.Errorf("opening store: %w", err)
			}

			es, err := store.Load(setName, environment)
			if err != nil {
				return fmt.Errorf("loading envset: %w", err)
			}

			switch action {
			case "add":
				if len(args) < 3 {
					return fmt.Errorf("'add' requires a tag name")
				}
				tag, err := envset.NewTag(args[2])
				if err != nil {
					return err
				}
				if err := envset.AddTag(es, tag); err != nil {
					return err
				}
				if err := store.Save(es); err != nil {
					return fmt.Errorf("saving envset: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Tag %q added to %q\n", args[2], setName)

			case "remove":
				if len(args) < 3 {
					return fmt.Errorf("'remove' requires a tag name")
				}
				if err := envset.RemoveTag(es, args[2]); err != nil {
					return err
				}
				if err := store.Save(es); err != nil {
					return fmt.Errorf("saving envset: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Tag %q removed from %q\n", args[2], setName)

			case "list":
				if outputJSON {
					enc := json.NewEncoder(cmd.OutOrStdout())
					return enc.Encode(es.Tags)
				}
				if len(es.Tags) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "(no tags)")
				} else {
					for _, t := range es.Tags {
						fmt.Fprintln(cmd.OutOrStdout(), t)
					}
				}

			default:
				return fmt.Errorf("unknown action %q: use add, remove, or list", action)
			}
			return nil
		},
	}

	tagCmd.Flags().StringVarP(&environment, "env", "e", "local", "Environment of the envset")
	tagCmd.Flags().BoolVar(&outputJSON, "json", false, "Output tags as JSON")

	rootCmd.AddCommand(tagCmd)
	_ = os.Stderr // suppress unused import
}
