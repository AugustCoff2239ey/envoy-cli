package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

var rollbackStacks = map[string]*envset.RollbackStack{}

func getRollbackStack(key string) *envset.RollbackStack {
	if _, ok := rollbackStacks[key]; !ok {
		rollbackStacks[key] = envset.NewRollbackStack(10)
	}
	return rollbackStacks[key]
}

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Manage rollback snapshots for an env set",
	}

	pushCmd := &cobra.Command{
		Use:   "push <name> <env>",
		Short: "Push current state onto the rollback stack",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, _ := cmd.Flags().GetString("message")
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("rollback push: %w", err)
			}
			key := args[0] + "/" + args[1]
			stack := getRollbackStack(key)
			if err := stack.Push(es, msg); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "pushed state for %s (stack depth: %d)\n", key, stack.Len())
			return nil
		},
	}
	pushCmd.Flags().StringP("message", "m", "", "message describing this rollback point")

	popCmd := &cobra.Command{
		Use:   "pop <name> <env>",
		Short: "Restore the most recent rollback state",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			es, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("rollback pop: %w", err)
			}
			key := args[0] + "/" + args[1]
			stack := getRollbackStack(key)
			if err := stack.Pop(es); err != nil {
				return err
			}
			if err := store.Save(es); err != nil {
				return fmt.Errorf("rollback pop save: %w", err)
			}
			fmt.Fprintf(os.Stdout, "restored previous state for %s\n", key)
			return nil
		},
	}

	peekCmd := &cobra.Command{
		Use:   "peek <name> <env>",
		Short: "Show the most recent rollback entry without restoring",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0] + "/" + args[1]
			stack := getRollbackStack(key)
			entry, err := stack.Peek()
			if err != nil {
				return err
			}
			out, _ := json.MarshalIndent(entry, "", "  ")
			fmt.Fprintln(os.Stdout, string(out))
			return nil
		},
	}

	rollbackCmd.AddCommand(pushCmd, popCmd, peekCmd)
	rootCmd.AddCommand(rollbackCmd)
}
