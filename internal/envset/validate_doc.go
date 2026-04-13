// Package envset provides primitives for creating, storing, diffing,
// syncing, exporting, importing, and validating environment variable sets.
//
// # Validation
//
// The Validate function inspects an [EnvSet] for common issues before the set
// is persisted or exported:
//
//   - Key names must conform to the POSIX convention: they must start with a
//     letter or underscore and contain only letters, digits, and underscores.
//   - A small set of well-known system keys (PATH, HOME, USER) are flagged as
//     reserved so that callers are warned before accidentally overwriting them.
//   - Values must not contain raw newline characters, which could break dotenv
//     and shell export formats.
//   - Values are capped at 4 096 characters to prevent accidental storage of
//     large blobs inside an env set.
//
// Example usage:
//
//	es, _ := envset.New("myapp", "production")
//	es.Vars["DB_URL"] = "postgres://localhost/prod"
//
//	if err := envset.Validate(es); err != nil {
//		log.Fatalf("invalid env set: %v", err)
//	}
//
// Individual keys can also be checked with [ValidateKey] before they are
// inserted into a set.
package envset
