/*
Package envset provides alias support for environment variable keys.

# Aliases

Aliases allow you to attach human-friendly or short names to existing keys
without duplicating values. An alias is stored as metadata on the key and
can be resolved back to the canonical key at any time.

Example usage:

	_ = envset.AddAlias(es, "DATABASE_HOST", "db-host")
	key, ok := envset.ResolveAlias(es, "db-host")
	// key == "DATABASE_HOST", ok == true

Alias names must start with a letter and may contain letters, digits,
underscores, or hyphens.
*/
package envset
