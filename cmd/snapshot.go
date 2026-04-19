package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var snapshotMessage string
	var restoreID string

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Take or restore a snapshot of an env set",
	}

	takeCmd := &cobra.Command{
		Use:   "take <name> <environment>",
		Short: "Take a snapshot of an env set and write it to stdout as JSON",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, env := args[0], args[1]
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("env set %q (%s) not found", name, env)
			}
			snap, err := envset.TakeSnapshot(es, snapshotMessage)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(snap)
		},
	}
	takeCmd.Flags().StringVarP(&snapshotMessage, "message", "m", "", "Optional snapshot message")

	restoreCmd := &cobra.Command{
		Use:   "restore <snapshot-file>",
		Short: "Restore an env set from a snapshot JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = restoreID
			snap, err := loadSnapshotFile(args[0])
			if err != nil {
				return err
			}
			es, err := envset.RestoreSnapshot(snap)
			if err != nil {
				return err
			}
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Restored %q (%s) from snapshot %s\n", es.Name, es.Environment, snap.ID)
			return nil
		},
	}
	restoreCmd.Flags().StringVar(&restoreID, "id", "", "Snapshot ID for reference")

	snapshotCmd.AddCommand(takeCmd, restoreCmd)
	rootCmd.AddCommand(snapshotCmd)
}

// loadSnapshotFile opens and decodes a snapshot JSON file from the given path.
func loadSnapshotFile(path string) (*envset.Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open snapshot file: %w", err)
	}
	defer f.Close()
	var snap envset.Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("invalid snapshot file: %w", err)
	}
	return &snap, nil
}
