package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, group, outputFmt string

	cmd := &cobra.Command{
		Use:   "index",
		Short: "Build and display a positional index of an env set",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("index: %w", err)
			}

			idx, err := envset.Index(es)
			if err != nil {
				return err
			}

			var entries []envset.IndexEntry
			if group != "" {
				entries = envset.IndexByGroup(idx, group)
			} else {
				for _, k := range envset.IndexKeys(idx) {
					entries = append(entries, idx[k])
				}
			}

			if outputFmt == "json" {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(entries)
			}

			for _, e := range entries {
				grpLabel := e.Group
				if grpLabel == "" {
					grpLabel = "-"
				}
				tagLabel := "-"
				if len(e.Tags) > 0 {
					tagLabel = fmt.Sprintf("%v", e.Tags)
				}
				fmt.Printf("[%3d] %-30s group=%-12s tags=%s\n",
					e.Position, e.Key, grpLabel, tagLabel)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	cmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Filter by group")
	cmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "Output format: text|json")
	_ = cmd.MarkFlagRequired("name")

	rootCmd.AddCommand(cmd)
}
