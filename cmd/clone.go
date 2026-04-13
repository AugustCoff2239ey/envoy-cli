package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/envset"
)

var (
	cloneNewName string
	cloneNewEnv  string
)

var cloneCmd = &cobra.Command{
	Use:   "clone <name> <environment>",
	Short: "Clone an existing env set into a new name or environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcName := args[0]
		srcEnv := args[1]

		store, err := envset.NewStore(storeDir)
		if err != nil {
			return fmt.Errorf("opening store: %w", err)
		}

		src, err := store.Load(srcName, srcEnv)
		if err != nil {
			return fmt.Errorf("loading source env set %q (%s): %w", srcName, srcEnv, err)
		}

		cloned, err := envset.Clone(src, envset.CloneOptions{
			NewName:        cloneNewName,
			NewEnvironment: cloneNewEnv,
		})
		if err != nil {
			return fmt.Errorf("cloning: %w", err)
		}

		if err := store.Save(cloned); err != nil {
			return fmt.Errorf("saving cloned env set: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Cloned %q (%s) → %q (%s) with %d variable(s).\n",
			srcName, srcEnv, cloned.Name, cloned.Environment, len(cloned.Values))
		return nil
	},
}

func init() {
	cloneCmd.Flags().StringVar(&cloneNewName, "new-name", "", "Name for the cloned env set (defaults to source name)")
	cloneCmd.Flags().StringVar(&cloneNewEnv, "new-env", "", "Environment for the cloned env set (defaults to source environment)")
	rootCmd.AddCommand(cloneCmd)
}
