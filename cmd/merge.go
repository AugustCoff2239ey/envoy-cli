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

		var strategy envset.MergeStrategy
		switch strings.ToLower(mergeStrategy) {
		case "ours":
			strategy = envset.MergeStrategyOurs
		case "theirs":
			strategy = envset.MergeStrategyTheirs
		case "error":
			strategy = envset.MergeStrategyError
		default:
			return fmt.Errorf("unknown strategy %q; use ours|theirs|error", mergeStrategy)
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

func init() {
	mergeCmd.Flags().StringVar(&mergeStrategy, "strategy", "ours", "conflict resolution strategy: ours|theirs|error")
	mergeCmd.Flags().StringVar(&mergeOutput, "output", "", "name for the resulting env set (defaults to base name)")
	rootCmd.AddCommand(mergeCmd)
}
