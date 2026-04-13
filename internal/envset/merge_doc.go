// Package envset provides types and utilities for managing named sets of
// environment variables scoped to a specific deployment environment.
//
// # Merge
//
// The Merge function combines two EnvSets (dst and src) into a new EnvSet.
// Keys present only in src are added to the result. Keys present in both sets
// are resolved according to the chosen MergeStrategy:
//
//   - MergeStrategyOurs   – the value from dst is kept (default-safe).
//   - MergeStrategyTheirs – the value from src overwrites dst.
//   - MergeStrategyError  – an error is returned immediately on the first
//     conflicting key, useful for strict CI pipelines.
//
// The original dst and src EnvSets are never mutated.
//
// Example:
//
//	result, err := envset.Merge(base, override, envset.MergeStrategyTheirs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("conflicts:", result.Conflicts)
package envset
