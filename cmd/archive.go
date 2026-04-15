package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

var globalArchive = envset.NewArchive()

func init() {
	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive and restore EnvSet snapshots",
	}

	addCmd := &cobra.Command{
		Use:   "add <name> <environment>",
		Short: "Archive the current state of an EnvSet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason, _ := cmd.Flags().GetString("reason")
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			entry, err := globalArchive.Add(es, reason)
			if err != nil {
				return fmt.Errorf("archive: %w", err)
			}
			fmt.Printf("Archived as %s\n", entry.ID)
			return nil
		},
	}
	addCmd.Flags().String("reason", "", "Reason for archiving")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all archive entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			asJSON, _ := cmd.Flags().GetBool("json")
			entries := globalArchive.List()
			if asJSON {
				return json.NewEncoder(os.Stdout).Encode(entries)
			}
			if len(entries) == 0 {
				fmt.Println("No archive entries.")
				return nil
			}
			for _, e := range entries {
				fmt.Printf("[%s] %s — %s\n", e.ArchivedAt.Format("2006-01-02 15:04:05"), e.ID, e.Reason)
			}
			return nil
		},
	}
	listCmd.Flags().Bool("json", false, "Output as JSON")

	restoreCmd := &cobra.Command{
		Use:   "restore <id> <name> <environment>",
		Short: "Restore an EnvSet from an archive entry",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			es, err := globalArchive.Restore(args[0], args[1], args[2])
			if err != nil {
				return fmt.Errorf("restore: %w", err)
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Printf("Restored %s/%s from archive %s\n", args[1], args[2], args[0])
			return nil
		},
	}

	archiveCmd.AddCommand(addCmd, listCmd, restoreCmd)
	rootCmd.AddCommand(archiveCmd)
}
