package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var format string

	statusCmd := &cobra.Command{
		Use:   "status <name> <environment>",
		Short: "Show live status of an env set",
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
			report, err := envset.Status(es)
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(report)
			}
			fmt.Printf("Name:        %s\n", report.Name)
			fmt.Printf("Environment: %s\n", report.Environment)
			fmt.Printf("Total Keys:  %d\n", report.TotalKeys)
			fmt.Printf("Locked:      %d\n", report.LockedKeys)
			fmt.Printf("Pinned:      %d\n", report.PinnedKeys)
			fmt.Printf("Protected:   %d\n", report.Protected)
			fmt.Printf("Expired:     %d\n", report.ExpiredKeys)
			fmt.Printf("Readonly:    %v\n", report.Readonly)
			fmt.Printf("Frozen:      %v\n", report.Frozen)
			fmt.Printf("Generated:   %s\n", report.GeneratedAt.Format("2006-01-02T15:04:05Z"))
			return nil
		},
	}

	statusCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(statusCmd)
}
