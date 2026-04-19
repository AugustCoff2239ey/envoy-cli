package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
se:<name> ort: "Display for an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeDir())
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("summary: %w", err)
			}

			r, err := envset.Summary(es)
			if err != nil {
				return err
			}

			outputJSON, _ := cmd.Flags().GetBool("json")
			if outputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(r)
			}

			printSummaryText(cmd, r)
			return nil
		},
	}

	summaryCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(summaryCmd)
}

func printSummaryText(cmd *cobra.Command, r *envset.SummaryReport) {
	fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", r.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "Environment: %s\n", r.Environment)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Keys:  %d\n", r.TotalKeys)
	fmt.Fprintf(cmd.OutOrStdout(), "Readonly:    %v\n", r.Readonly)
	fmt.Fprintf(cmd.OutOrStdout(), "Frozen:      %v\n", r.Frozen)
	if len(r.LockedKeys) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Locked:      %s\n", strings.Join(r.LockedKeys, ", "))
	}
	if len(r.Protected) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Protected:   %s\n", strings.Join(r.Protected, ", "))
	}
	if len(r.Pinned) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Pinned:      %s\n", strings.Join(r.Pinned, ", "))
	}
	if len(r.Expired) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Expired:     %s\n", strings.Join(r.Expired, ", "))
	}
	if len(r.TagList) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Tags:        %s\n", strings.Join(r.TagList, ", "))
	}
}
