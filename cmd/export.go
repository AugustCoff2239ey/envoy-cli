package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

var exportCmd = &cobra.Command{
	Use:   "export <name> <environment>",
	Short: "Export an env set to a given format",
	Long:  `Export an environment variable set to dotenv, shell, or JSON format.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runExport,
}

var exportFormat string
var exportOutput string

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, shell, json")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	name := args[0]
	env := args[1]

	store, err := envset.NewStore(defaultStorePath())
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}

	set, err := store.Load(name, env)
	if err != nil {
		return fmt.Errorf("env set %q (%s) not found: %w", name, env, err)
	}

	output, err := envset.Export(set, exportFormat)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	if exportOutput == "" {
		fmt.Print(output)
		return nil
	}

	if err := os.WriteFile(exportOutput, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Exported %q (%s) to %s\n", name, env, exportOutput)
	return nil
}
