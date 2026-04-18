package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

var globalCheckpointStore = envset.NewCheckpointStore()

func init() {
	checkpointCmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Manage named checkpoints for an env set",
	}

	var cpEnv string

	saveCmd := &cobra.Command{
		Use:   "save <name> <set>",
		Short: "Save a checkpoint of an env set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cpName, setName := args[0], args[1]
			st := newStore()
			e, err := st.Load(setName, cpEnv)
			if err != nil {
				return fmt.Errorf("set not found: %w", err)
			}
			if err := globalCheckpointStore.Save(cpName, e); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q saved for %s/%s\n", cpName, setName, cpEnv)
			return nil
		},
	}
	saveCmd.Flags().StringVar(&cpEnv, "env", "local", "environment of the set")

	restoreCmd := &cobra.Command{
		Use:   "restore <name> <set>",
		Short: "Restore a checkpoint onto an env set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cpName, setName := args[0], args[1]
			st := newStore()
			e, err := st.Load(setName, cpEnv)
			if err != nil {
				return fmt.Errorf("set not found: %w", err)
			}
			if err := globalCheckpointStore.Restore(cpName, e); err != nil {
				return err
			}
			if err := st.Save(e); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q restored onto %s/%s\n", cpName, setName, cpEnv)
			return nil
		},
	}
	restoreCmd.Flags().StringVar(&cpEnv, "env", "local", "environment of the set")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved checkpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			names := globalCheckpointStore.List()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No checkpoints saved.")
				return nil
			}
			outFmt, _ := cmd.Flags().GetString("output")
			if outFmt == "json" {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(names)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			}
			return nil
		},
	}
	listCmd.Flags().String("output", "text", "output format: text|json")

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named checkpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := globalCheckpointStore.Delete(args[0]); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q deleted\n", args[0])
			return nil
		},
	}

	checkpointCmd.AddCommand(saveCmd, restoreCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(checkpointCmd)
}
