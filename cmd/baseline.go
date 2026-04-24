package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, label, format string

	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Capture or check drift against a baseline snapshot",
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Capture current state as a named baseline",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}
			b, err := envset.SetBaseline(e, label)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Baseline %q captured at %s (%d keys)\n",
				b.Label, b.CreatedAt.Format("2006-01-02T15:04:05Z"), len(b.Vars))
			return nil
		},
	}
	setCmd.Flags().StringVar(&name, "name", "", "EnvSet name (required)")
	setCmd.Flags().StringVar(&env, "env", "local", "Environment")
	setCmd.Flags().StringVar(&label, "label", "", "Baseline label (required)")
	_ = setCmd.MarkFlagRequired("name")
	_ = setCmd.MarkFlagRequired("label")

	driftCmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect drift between current state and a baseline",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return err
			}
			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("envset %q (%s) not found", name, env)
			}
			b, err := envset.SetBaseline(e, label)
			if err != nil {
				return err
			}
			result, err := envset.DetectDrift(e, b)
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(result)
			}
			if !result.HasDrift() {
				fmt.Println("No drift detected.")
				return nil
			}
			for k, v := range result.Added {
				fmt.Printf("+ %s=%s\n", k, v)
			}
			for k, v := range result.Removed {
				fmt.Printf("- %s=%s\n", k, v)
			}
			for k, pair := range result.Changed {
				fmt.Printf("~ %s: %s -> %s\n", k, pair[0], pair[1])
			}
			return nil
		},
	}
	driftCmd.Flags().StringVar(&name, "name", "", "EnvSet name (required)")
	driftCmd.Flags().StringVar(&env, "env", "local", "Environment")
	driftCmd.Flags().StringVar(&label, "label", "current", "Baseline label")
	driftCmd.Flags().StringVar(&format, "format", "text", "Output format: text|json")
	_ = driftCmd.MarkFlagRequired("name")

	baselineCmd.AddCommand(setCmd, driftCmd)
	rootCmd.AddCommand(baselineCmd)
}
