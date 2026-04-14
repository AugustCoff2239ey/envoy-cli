package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envset"
)

func init() {
	var name, env, schemaFile string
	var outputJSON bool

	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate an envset against a JSON schema file",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := envset.NewStore(defaultStorePath())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			es, err := store.Load(name, env)
			if err != nil {
				return fmt.Errorf("load envset %q (%s): %w", name, env, err)
			}

			data, err := os.ReadFile(schemaFile)
			if err != nil {
				return fmt.Errorf("read schema file: %w", err)
			}

			var raw map[string]struct {
				Type     string `json:"type"`
				Required bool   `json:"required"`
				Pattern  string `json:"pattern"`
			}
			if err := json.Unmarshal(data, &raw); err != nil {
				return fmt.Errorf("parse schema: %w", err)
			}

			schema := make(envset.Schema, len(raw))
			for k, v := range raw {
				schema[k] = envset.SchemaField{
					Type:     envset.SchemaFieldType(v.Type),
					Required: v.Required,
					Pattern:  v.Pattern,
				}
			}

			violations, err := envset.ValidateSchema(es, schema)
			if err != nil {
				return fmt.Errorf("validate schema: %w", err)
			}

			if outputJSON {
				return json.NewEncoder(os.Stdout).Encode(violations)
			}

			if len(violations) == 0 {
				fmt.Println("✓ All schema constraints satisfied.")
				return nil
			}

			fmt.Printf("Schema violations (%d):\n", len(violations))
			for _, v := range violations {
				fmt.Printf("  [%s] %s\n", v.Key, v.Message)
			}
			return fmt.Errorf("schema validation failed")
		},
	}

	schemaCmd.Flags().StringVarP(&name, "name", "n", "", "EnvSet name (required)")
	schemaCmd.Flags().StringVarP(&env, "env", "e", "local", "Environment (local/staging/production)")
	schemaCmd.Flags().StringVarP(&schemaFile, "file", "f", "", "Path to JSON schema file (required)")
	schemaCmd.Flags().BoolVar(&outputJSON, "json", false, "Output violations as JSON")
	_ = schemaCmd.MarkFlagRequired("name")
	_ = schemaCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(schemaCmd)
}
