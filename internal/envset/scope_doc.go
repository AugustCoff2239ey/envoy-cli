/*
Package envset provides scope management for grouping environment variable keys
into named namespaces within an EnvSet.

# Scopes

A scope is a named collection of key references inside an EnvSet. Scopes allow
operators to logically partition variables — for example separating database
credentials from API settings — without duplicating values.

Scopes are stored in the EnvSet.Scopes map (map[string][]string). Keys in a
scope must already exist in the EnvSet at creation time.

# Usage

	err := envset.CreateScope(es, "database", []string{"DB_HOST", "DB_PORT", "DB_NAME"})

	keys, err := envset.GetScope(es, "database")

	vars, err := envset.ScopeVars(es, "database")

	names := envset.ListScopes(es)

	err = envset.DeleteScope(es, "database")

# Scope Name Rules

Scope names must begin with a letter and contain only letters, digits,
underscores, or hyphens. Empty names and numeric-leading names are rejected.
*/
package envset
