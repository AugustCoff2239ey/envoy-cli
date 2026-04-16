// Package envset provides the Annotate feature for attaching human-readable
// notes to individual environment variable keys within an EnvSet.
//
// # Overview
//
// Annotations are short, single-line notes stored alongside an EnvSet's
// variables. They are useful for documenting the purpose of a key, noting
// rotation schedules, or flagging keys that require special handling.
//
// # Usage
//
//	// Attach a note
//	err := envset.Annotate(e, "API_KEY", "Rotated every 90 days")
//
//	// Retrieve a note
//	note, ok := envset.GetAnnotation(e, "API_KEY")
//
//	// List all annotations
//	annotations := envset.ListAnnotations(e)
//
//	// Remove a note
//	err = envset.RemoveAnnotation(e, "API_KEY")
//
// # Constraints
//
//   - Notes must not contain newline characters.
//   - The target key must already exist in the EnvSet.
//   - The EnvSet must not be nil.
package envset
