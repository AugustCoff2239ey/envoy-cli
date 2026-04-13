package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var newName string
	var newEnv string

	renameCmd := &cobra.Command{
		Use:   "rename <name> <environment>",
		Short: "Rename an envset's name or environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			environment := args[1]

			if newName == "" && newEnv == "" {
				return fmt.Errorf("at least one of --new-name or --new-env must be provided")
			}

			store, err := envset.NewStore(storeDir())
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}

			src, err := store.Load(name, environment)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found: %w", name, environment, err)
			}

			opts := envset.RenameOptions{
				NewName:        newName,
				NewEnvironment: newEnv,
			}

			renamed, err := envset.Rename(store, src, opts)
			if err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Renamed to %q (%s)\n", renamed.Name, renamed.Environment)
			return nil
		},
	}

	renameCmd.Flags().StringVar(&newName, "new-name", "", "New name for the envset")
	renameCmd.Flags().StringVar(&newEnv, "new-env", "", "New environment for the envset")

	rootCmd.AddCommand(renameCmd)
}
