/*
Package envset provides group management for EnvSet variables.

# Groups

A group is a named subset of keys within an EnvSet, stored in the EnvSet's
metadata under the "group:<name>" key. Groups allow users to logically
organize related environment variables (e.g., "database", "auth", "cache")
without affecting the underlying key-value pairs.

Group names must contain only letters, digits, hyphens, and underscores.
All keys referenced in a group must already exist in the EnvSet at creation
time.

# Example

	es, _ := envset.New("myapp", "production")
	_ = es.Set("DB_HOST", "db.example.com")
	_ = es.Set("DB_PORT", "5432")

	_ = envset.CreateGroup(es, "database", []string{"DB_HOST", "DB_PORT"})

	g, _ := envset.GetGroup(es, "database")
	fmt.Println(g.Keys) // [DB_HOST DB_PORT]

	groups := envset.ListGroups(es)
	for _, grp := range groups {
		fmt.Printf("Group %s: %v\n", grp.Name, grp.Keys)
	}
*/
package envset
