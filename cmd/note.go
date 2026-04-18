package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, author string

	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage notes attached to an env set",
	}

	addCmd := &cobra.Command{
		Use:   "add <text>",
		Short: "Add a note to an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("env set not found: %w", err)
			}
			if err := envset.AddNote(es, args[0], author); err != nil {
				return err
			}
			return st.Save(es)
		},
	}
	addCmd.Flags().StringVar(&author, "author", "", "Author of the note")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List notes on an env set",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("env set not found: %w", err)
			}
			notes := envset.ListNotes(es)
			if len(notes) == 0 {
				fmt.Fprintln(os.Stdout, "no notes")
				return nil
			}
			for i, n := range notes {
				fmt.Fprintf(os.Stdout, "[%d] %s (%s) — %s\n", i+1, n.Text, n.Author, n.CreatedAt.Format("2006-01-02 15:04"))
			}
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all notes from an env set",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := envset.NewStore(storeFile)
			if err != nil {
				return err
			}
			es, err := st.Load(name, env)
			if err != nil {
				return fmt.Errorf("env set not found: %w", err)
			}
			if err := envset.ClearNotes(es); err != nil {
				return err
			}
			return st.Save(es)
		},
	}

	for _, sub := range []*cobra.Command{addCmd, listCmd, clearCmd} {
		sub.Flags().StringVar(&name, "name", "", "Env set name (required)")
		sub.Flags().StringVar(&env, "env", "local", "Environment")
		_ = sub.MarkFlagRequired("name")
		noteCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(noteCmd)
}
