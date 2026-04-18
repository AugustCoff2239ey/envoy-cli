package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var note string

	draftCmd := &cobra.Command{
		Use:   "draft",
		Short: "Manage draft envsets",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name> <env>",
		Short: "Save a draft of an envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			e, err := store.Load(args[0], args[1])
			if err != nil {
				return fmt.Errorf("envset not found: %w", err)
			}
			d, err := envset.SaveDraft(e, note)
			if err != nil {
				return err
			}
			if err := store.Save(d); err != nil {
				return err
			}
			fmt.Printf("Draft saved as %q\n", d.Name)
			return nil
		},
	}
	saveCmd.Flags().StringVar(&note, "note", "", "Optional note for the draft")

	promoteCmd := &cobra.Command{
		Use:   "promote <name> <env>",
		Short: "Promote a draft to a regular envset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			draftName := "draft:" + args[0]
			d, err := store.Load(draftName, args[1])
			if err != nil {
				return fmt.Errorf("draft not found: %w", err)
			}
			p, err := envset.PromoteDraft(d)
			if err != nil {
				return err
			}
			if err := store.Save(p); err != nil {
				return err
			}
			_ = store.Delete(draftName, args[1])
			fmt.Printf("Draft promoted to %q\n", p.Name)
			return nil
		},
	}

	inspectCmd := &cobra.Command{
		Use:   "inspect <name> <env>",
		Short: "Inspect a draft envset as JSON",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(storeFile())
			if err != nil {
				return err
			}
			d, err := store.Load("draft:"+args[0], args[1])
			if err != nil {
				return fmt.Errorf("draft not found: %w", err)
			}
			return json.NewEncoder(os.Stdout).Encode(d)
		},
	}

	draftCmd.AddCommand(saveCmd, promoteCmd, inspectCmd)
	rootCmd.AddCommand(draftCmd)
}
