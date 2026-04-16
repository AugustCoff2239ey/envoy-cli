package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/envset"
)

var (
	mergeStrategy string
	mergeOutput   string
)

var mergeCmd = &cobra.Command{
	Use:   "merge <base-name> <source-name>",
	Short: "Merge two env sets into one",
	Long: `Merge combines the keys from the source env set into the base env set.
Conflicts are resolved according to --strategy (ours|theirs|error).`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseName, srcName := args[0], args[1]

		if baseName == srcName {
			return fmt.Errorf("base and source env sets must be different")
		}

		st, err := envset.NewStore(storeDir())
		if err != nil {
			return err
		}

		base, err := st.Load(baseName)
		if err != nil {
			return fmt.Errorf("loading base %q: %w", baseName, err)
		}
		src, err := st.Load(srcName)
		if err != nil {
			return fmt.Errorf("loading source %q: %w", srcName, err)
		}

		strategy, err := parseMergeStrategy(mergeStrategy)
		if err != nil {
			return err
		}

		result, err := envset.Merge(base, src, strategy)
		if err != nil {
			return err
		}

		if len(result.Conflicts) > 0 {
			fmt.Fprintf(os.Stderr, "conflicts resolved (%s): %s\n",
				mergeStrategy, strings.Join(result.Conflicts, ", "))
		}

		if mergeOutput != "" {
			result.Merged.Name = mergeOutput
		}

		if err := st.Save(result.Merged); err != nil {
			return fmt.Errorf("saving merged set: %w", err)
		}

		fmt.Printf("merged %q + %q → %q\n", baseName, srcName, result.Merged.Name)
		return nil
	},
}

// parseMergeStrategy converts a strategy string to the corresponding
// envset.MergeStrategy constant, returning an error for unknown values.
func parseMergeStrategy(s string) (envset.MergeStrategy, error) {
	switch strings.ToLower(s) {
	case "ours":
		return envset.MergeStrategyOurs, nil
	case "theirs":
		return envset.MergeStrategyTheirs, nil
	case "error":
		return envset.MergeStrategyError, nil
	default:
		return 0, fmt.Errorf("unknown strategy %q; use ours|theirs|error", s)
	}
}

func init() {
	mergeCmd.Flags().StringVar(&mergeStrategy, "strategy", "ours", "conflict resolution strategy: ours|theirs|error")
	mergeCmd.Flags().StringVar(&mergeOutput, "output", "", "name for the resulting env set (defaults to base name)")
	rootCmd.AddCommand(mergeCmd)
}
