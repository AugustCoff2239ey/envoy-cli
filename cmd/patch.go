package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var setFlags []string
	var deleteFlags []string
	var fromDiff string

	patchCmd := &cobra.Command{
		Use:   "patch <name> <environment>",
		Short: "Apply patch operations to an env set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, env := args[0], args[1]
			store := envset.NewStore(defaultStorePath())

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("patch: load %q/%q: %w", name, env, err)
			}

			var entries []envset.PatchEntry

			if fromDiff != "" {
				data, err := os.ReadFile(fromDiff)
				if err != nil {
					return fmt.Errorf("patch: read diff file: %w", err)
				}
				var d envset.DiffResult
				if err := json.Unmarshal(data, &d); err != nil {
					return fmt.Errorf("patch: parse diff file: %w", err)
				}
				entries = envset.PatchFromDiff(d)
			}

			for _, s := range setFlags {
				k, v, ok := splitKV(s)
				if !ok {
					return fmt.Errorf("patch: invalid --set value %q (expected KEY=VALUE)", s)
				}
				entries = append(entries, envset.PatchEntry{Op: envset.PatchOpSet, Key: k, Value: v})
			}
			for _, k := range deleteFlags {
				entries = append(entries, envset.PatchEntry{Op: envset.PatchOpDelete, Key: k})
			}

			if len(entries) == 0 {
				return fmt.Errorf("patch: no operations provided")
			}

			if err := envset.Patch(es, entries); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("patch: save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "patched %d operation(s) on %q/%q\n", len(entries), name, env)
			return nil
		},
	}

	patchCmd.Flags().StringArrayVar(&setFlags, "set", nil, "Set a key (KEY=VALUE)")
	patchCmd.Flags().StringArrayVar(&deleteFlags, "delete", nil, "Delete a key")
	patchCmd.Flags().StringVar(&fromDiff, "from-diff", "", "Apply patch from a JSON diff file")

	rootCmd.AddCommand(patchCmd)
}
