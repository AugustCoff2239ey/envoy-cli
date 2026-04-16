package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, format string

	annotateCmd := &cobra.Command{
		Use:   "annotate <key> [note]",
		Short: "Attach or remove notes on environment variable keys",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := envset.NewStore(storeDir())
			e, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("annotate: %w", err)
			}

			remove, _ := cmd.Flags().GetBool("remove")
			list, _ := cmd.Flags().GetBool("list")

			switch {
			case list:
				annotations := envset.ListAnnotations(e)
				if format == "json" {
					return json.NewEncoder(os.Stdout).Encode(annotations)
				}
				for _, a := range annotations {
					fmt.Printf("%s: %s\n", a.Key, a.Note)
				}
			case remove:
				if err := envset.RemoveAnnotation(e, args[0]); err != nil {
					return err
				}
				if err := store.Save(e); err != nil {
					return err
				}
				fmt.Printf("annotation removed from %q\n", args[0])
			default:
				if len(args) < 2 {
					return fmt.Errorf("annotate: note required when not using --remove or --list")
				}
				if err := envset.Annotate(e, args[0], args[1]); err != nil {
					return err
				}
				if err := store.Save(e); err != nil {
					return err
				}
				fmt.Printf("annotated %q\n", args[0])
			}
			return nil
		},
	}

	annotateCmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	annotateCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment")
	annotateCmd.Flags().StringVar(&format, "format", "text", "Output format: text|json")
	annotateCmd.Flags().Bool("remove", false, "Remove annotation from key")
	annotateCmd.Flags().Bool("list", false, "List all annotations")
	_ = annotateCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(annotateCmd)
}
